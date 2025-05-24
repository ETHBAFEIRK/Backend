import requests
import subprocess

def fetch_rates(url="http://localhost:8080/rates"):
    resp = requests.get(url)
    resp.raise_for_status()
    return resp.json()

def build_mermaid_graph(rates):
    # --- Get the set of tokens with icons from the backend ---
    # This list must match the backend's tokenIcons keys.
    token_icons = {
        "ETH", "WETH", "stETH", "wstETH", "ezETH", "pzETH", "STONE", "xPufETH", "mstETH", "weETH", "egETH",
        "inwstETH", "rsETH", "LsETH", "USDC", "USDT", "USDe", "FBTC", "LBTC", "mBTC", "pumpBTC", "mswETH",
        "mwBETH", "mETH", "rstETH", "steakLRT", "Re7LRT", "amphrETH", "rswETH", "swETH", "weETHs"
    }

    # Build the graph as adjacency list and reverse adjacency for pruning
    nodes = set()
    edge_map = {}
    target_apy = {}

    for rate in rates:
        from_token = rate.get("input_symbol")
        to_token = rate.get("output_token")
        kind = rate.get("output_kind", "")
        apy = rate.get("apy")
        nodes.add(from_token)
        nodes.add(to_token)
        edge_map.setdefault((from_token, to_token), []).append((kind, apy))
        if to_token not in target_apy or apy > target_apy[to_token]:
            target_apy[to_token] = apy

    # --- Prune nodes and edges not leading to a token with an icon ---
    # 1. Build adjacency list
    adj = {}
    rev_adj = {}
    for (from_token, to_token) in edge_map:
        adj.setdefault(from_token, set()).add(to_token)
        rev_adj.setdefault(to_token, set()).add(from_token)

    # 2. Find all nodes that can reach a token with an icon (reverse BFS)
    reachable = set(token_icons)
    queue = list(token_icons)
    while queue:
        curr = queue.pop()
        for prev in rev_adj.get(curr, []):
            if prev not in reachable:
                reachable.add(prev)
                queue.append(prev)

    # 3. Only keep nodes and edges where the output token is in reachable set and is in token_icons
    pruned_nodes = set()
    pruned_edges = []
    for (from_token, to_token), kind_apys in edge_map.items():
        if to_token not in reachable:
            continue
        if to_token not in token_icons:
            continue
        pruned_nodes.add(from_token)
        pruned_nodes.add(to_token)
        # Prefer "stake"/"restake" over "swap"
        preferred = None
        for kind, apy in kind_apys:
            if kind in ("stake", "restake"):
                preferred = kind
                break
        if preferred:
            label = preferred
        else:
            label = kind_apys[0][0]
        label = label.replace('"', '\\"') if label else ""
        pruned_edges.append((from_token, to_token, label))

    # 4. Only keep nodes that are in pruned_edges
    mermaid = ["graph TD"]
    for node in sorted(pruned_nodes):
        if node in target_apy:
            apy_val = target_apy[node]
            mermaid.append(f'    {node}["{node} ({apy_val:.2f}%)"]')
        else:
            mermaid.append(f'    {node}["{node}"]')
    for from_token, to_token, label in pruned_edges:
        if label:
            mermaid.append(f'    {from_token} -->|{label}| {to_token}')
        else:
            mermaid.append(f'    {from_token} --> {to_token}')
    return "\n".join(mermaid)

def render_mermaid(input_path: str, output_path: str):
    subprocess.run([
        "mmdc",
        "-i", input_path,
        "-o", output_path,
        "-b", "white",
        "--width", "2048"
    ], check=True)

def show_image(path: str):
    subprocess.run(["open", path], check=True)

def main():
    rates = fetch_rates()
    mermaid_graph = build_mermaid_graph(rates)
    mmd_path = "tokens_graph.mmd"
    png_path = "tokens_graph.png"
    with open(mmd_path, "w") as f:
        f.write(mermaid_graph)
    print("Mermaid graph written to tokens_graph.mmd")
    render_mermaid(mmd_path, png_path)
    print("PNG graph written to tokens_graph.png")
    show_image(png_path)

if __name__ == "__main__":
    main()
