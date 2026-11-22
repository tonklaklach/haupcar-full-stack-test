import type { Metadata } from "next";
import "./globals.css";
export const metadata: Metadata = {
  title: "HAUP | จัดการข้อมูลรถยนต์",
  description: "ทะเบียนรถยนต์ของบริษัท",
};
export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="th">
      <body>{children}</body>
    </html>
  );
}
