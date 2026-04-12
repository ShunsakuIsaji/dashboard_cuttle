import requests
from pathlib import Path

def fetch_file(url: str, save_path: Path) -> Path:
    response = requests.get(url)
    response.raise_for_status()  # Check if the request was successful

    save_path.parent.mkdir(parents=True, exist_ok=True)  # Create directories if they don't exist
    with open(save_path, 'wb') as file:
        file.write(response.content)

    return save_path