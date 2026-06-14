from fetch_file import fetch_file
from extract_xlsx import extract_data_from_xlsx
from gcs import fetch_gcs_file_to_df, upload_df_to_gcs
from dotenv import load_dotenv
from pathlib import Path
from typing import Any
import pandas as pd
import yaml
import os

CONFIG_PATH = Path("item_meta.yaml")
RAW_DIR = Path("data/raw")

def load_config(path: Path) -> dict[str, Any]:
    with path.open("r", encoding="utf-8") as f:
        return yaml.safe_load(f)

def fetch_data_from_endpoint(config: dict[str, Any], source_name: str) -> Path:
    source_config = config["sources"][source_name]
    url = source_config["url"]
    filename = source_config["filename"]

    save_path = RAW_DIR / f"{filename}"

    print(f"[INFO] Fetching data for source: {source_name}")
    print(f"[INFO] URL: {url}")

    fetched_path = fetch_file(url, save_path)
    print(f"[INFO] Data fetched and saved to: {fetched_path}")

    return fetched_path

def renew_gcs_data(df: pd.DataFrame, gcp_path: str) -> None:
    df_old = fetch_gcs_file_to_df(gcp_path)
    df_combined = pd.concat([df_old, df], ignore_index=True)
    df_combined = df_combined.drop_duplicates(subset=["date", "item"], keep = "last")
    print(f"[INFO] Uploading processed data to GCS at: {gcp_path}")
    upload_df_to_gcs(df_combined, gcp_path)
    print(f"[INFO] Data uploaded to GCS successfully.")

def main() -> None:
    if os.path.exists(".env"):
        load_dotenv()
    else:
        print("[INFO] .env file not found. Use environment variables for configuration.")
    
    config = load_config(CONFIG_PATH)
    df = pd.DataFrame()

    for source_name, source_config in config["sources"].items():
        try:
            fetched_path = fetch_data_from_endpoint(config, source_name)
            df_source = extract_data_from_xlsx(
                file_path=fetched_path,
                sheet_name=source_config["sheet"],
                skip_rows=source_config["skip_rows"],
                date_column=source_config["date_column"],
                items=source_config["items"]
            )
            df = pd.concat([df, df_source], ignore_index=True)
            df_source.to_csv(f"data/processed/{source_name}.csv", index=False)
            print(f"[INFO] Data for source {source_name} processed and saved to CSV.")

        except Exception as e:
            print(f"[ERROR] Failed to process source {source_name}: {e}")
    
    df = df.sort_values(["date", "item"]).reset_index(drop=True)

    renew_gcs_data(df, "all_data.csv")
    
    df.to_csv("data/all_data.csv", index=False)
    print("[INFO] All data processed and saved to data/all_data.csv")


if __name__ == "__main__":
    main()

