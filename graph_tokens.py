import requests
import subprocess

def fetch_rates(url="http://localhost:8080/rates"):
    resp = requests.get(url)
    resp.raise_for_status()
    return resp.json()

def build_mermaid_graph(rates):
    nodes = set()
    # For each (from_token, to_token), collect all rates and their kinds
    edge_map = {}
    # Map to store the highest APY for each target node (output token)
    target_apy = {}

    for rate in rates:
        from_token = rate.get("input_symbol")
        to_token = rate.get("output_token")
        kind = rate.get("output_kind", "")
        apy = rate.get("apy")
        nodes.add(from_token)
        nodes.add(to_token)
        # Store all rates for each (from, to) pair
        edge_map.setdefault((from_token, to_token), []).append((kind, apy))
        # Store the highest APY for each output token
        if to_token not in target_apy or apy > target_apy[to_token]:
            target_apy[to_token] = apy

    # Now, for each (from, to), prefer "stake"/"restake" over "swap"
    edges = []
    for (from_token, to_token), kind_apys in edge_map.items():
        # Find if any "stake" or "restake" exists
        preferred = None
        preferred_apy = None
        for kind, apy in kind_apys:
            if kind in ("stake", "restake"):
                preferred = kind
                preferred_apy = apy
                break
        if preferred:
            label = preferred
        else:
            # If no stake/restake, use the first kind (likely "swap")
            label = kind_apys[0][0]
        # Escape quotes in label
        label = label.replace('"', '\\"') if label else ""
        edges.append((from_token, to_token, label))

    mermaid = ["graph TD"]
    for node in sorted(nodes):
        if node in target_apy:
            apy_val = target_apy[node]
            # Format APY to 2 decimal places, show as e.g. "wstETH (4.12%)"
            mermaid.append(f'    {node}["{node} ({apy_val:.2f}%)"]')
        else:
            mermaid.append(f'    {node}["{node}"]')
    for from_token, to_token, label in edges:
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
