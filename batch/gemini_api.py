from pathlib import Path
from google import genai
from google.genai import types
from typing import Any
import os
import mimetypes


def build_prompt(source_name: str, source_config: dict[str, Any]) -> str:
    item_lines = []
    for item_key, meta in source_config["items"].items():
        item_lines.append(
            f"- {item_key}: {meta['description']} "
            f"(category={meta['category']}, unit={meta['unit']})"
        )

    item_text = "\n".join(item_lines)

    return f"""
あなたは表データ抽出アシスタントです。
以下の資料から、月次データのみを抽出してください。

対象ソース:
{source_name}

使用可能な item は以下のみです。
新しい item を作らないでください。不明なものは出力しないでください。

{item_text}

出力ルール:
- JSONのみを返す
- schema は以下
{{
  "records": [
    {{
      "date": "YYYY-MM",
      "item": "string",
      "value": number
    }}
  ]
}}
- 年平均は除外
- 月次のみ抽出
- value は数値のみ
- item は上の候補からのみ選択
""".strip()


def detect_mime_type(path: Path) -> str:
    mime_type, _ = mimetypes.guess_type(path.name)
    if mime_type:
        return mime_type

    suffix = path.suffix.lower()
    if suffix == ".xlsx":
        return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
    if suffix == ".pdf":
        return "application/pdf"
    return "application/octet-stream"


def send_to_gemini(file_path: Path, prompt: str) -> str:
    api_key = os.environ.get("GEMINI_API_KEY")
    if not api_key:
        raise RuntimeError("GEMINI_API_KEY is not set")

    client = genai.Client(api_key=api_key)

    mime_type = detect_mime_type(file_path)

    uploaded_file = client.files.upload(
        file=str(file_path),
        config=types.UploadFileConfig(mime_type=mime_type),
    )

    response = client.models.generate_content(
        model="gemini-2.5-flash",
        contents=[
            uploaded_file,
            prompt,
        ],
        config=types.GenerateContentConfig(
            response_mime_type="application/json",
        ),
    )

    if not response.text:
        raise RuntimeError("Gemini returned empty response")

    return response.text
