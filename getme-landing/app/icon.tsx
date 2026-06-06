import { ImageResponse } from "next/og";
import fs from "fs";
import path from "path";

export const size = {
  width: 32,
  height: 32,
};
export const contentType = "image/png";

export default function Icon() {
  const logoData = fs.readFileSync(
    path.join(process.cwd(), "public", "icon.png"),
  );
  const base64Image = logoData.toString("base64");
  const dataUrl = `data:image/png;base64,${base64Image}`;

  return new ImageResponse(
    <div tw="w-full h-full flex items-center justify-center rounded-lg overflow-hidden bg-transparent">
      <img src={dataUrl} width="100%" height="100%" alt="Icon" />
    </div>,
    {
      ...size,
    },
  );
}
