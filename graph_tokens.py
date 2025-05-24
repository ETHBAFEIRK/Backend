import requests

def fetch_rates(url="http://localhost:8080/rates"):
    resp = requests.get(url)
    resp.raise_for_status()
    return resp.json()

def build_mermaid_graph(rates):
    nodes = set()
    edges = []
    for rate in rates:
        from_token = rate.get("input_symbol")
        to_token = rate.get("output_token")
        kind = rate.get("output_kind", "")
        nodes.add(from_token)
        nodes.add(to_token)
        label = kind if kind else ""
        # Escape quotes in label
        label = label.replace('"', '\\"')
        edges.append((from_token, to_token, label))
    mermaid = ["graph TD"]
    for node in sorted(nodes):
        mermaid.append(f'    {node}["{node}"]')
    for from_token, to_token, label in edges:
        if label:
            mermaid.append(f'    {from_token} -->|{label}| {to_token}')
        else:
            mermaid.append(f'    {from_token} --> {to_token}')
    return "\n".join(mermaid)

def main():
    import mermaid
    rates = fetch_rates()
    mermaid_graph = build_mermaid_graph(rates)
    mmd_path = "tokens_graph.mmd"
    svg_path = "tokens_graph.svg"
    with open(mmd_path, "w") as f:
        f.write(mermaid_graph)
    print("Mermaid graph written to tokens_graph.mmd")
    # Render to SVG using mermaid-py
    svg = mermaid.render(mermaid_graph, output_format="svg")
    with open(svg_path, "w") as f:
        f.write(svg)
    print("SVG graph written to tokens_graph.svg")

if __name__ == "__main__":
    main()
