"use client";

import { useEffect, useRef, useState } from "react";
import { carsAPI, type Car } from "../lib/cars";
import { registrationSearchKey } from "../lib/registration";
import { CarDialog, type DialogAction } from "./car-dialog";

const dateFormat = new Intl.DateTimeFormat("th-TH", {
  dateStyle: "medium",
  timeZone: "Asia/Bangkok",
});

function searchKey(text: string) {
  return text.toLowerCase().replace(/\s/g, "");
}

export default function CarsPage() {
  const [cars, setCars] = useState<Car[]>([]);
  const [query, setQuery] = useState("");
  const [load, setLoad] = useState({ pending: true, error: "" });
  const [attempt, setAttempt] = useState(0);
  const [action, setAction] = useState<DialogAction | null>(null);
  const [notice, setNotice] = useState("");
  const addButton = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    let active = true;
    carsAPI.list().then(
      (data) => {
        if (!active) return;
        setCars(data);
        setLoad({ pending: false, error: "" });
      },
      (error: Error) => {
        if (active) setLoad({ pending: false, error: error.message });
      },
    );
    return () => {
      active = false;
    };
  }, [attempt]);

  const needle = searchKey(query);
  const plateNeedle = registrationSearchKey(query);
  const visible = cars.filter(
    (car) =>
      (!!plateNeedle &&
        registrationSearchKey(car.registration).includes(plateNeedle)) ||
      [car.brand, car.model, car.notes].some((value) =>
        searchKey(value).includes(needle),
      ),
  );

  function saved(car: Car) {
    if (action?.kind === "add") {
      setCars((previous) => [car, ...previous]);
      setQuery("");
      setNotice("เพิ่มข้อมูลรถยนต์แล้ว");
    } else {
      setCars((previous) =>
        previous.map((item) => (item.id === car.id ? car : item)),
      );
      setNotice("บันทึกการแก้ไขแล้ว");
    }
  }

  function removed(id: number) {
    setCars((previous) => previous.filter((car) => car.id !== id));
    setNotice("ลบข้อมูลรถยนต์แล้ว");
    requestAnimationFrame(() => addButton.current?.focus());
  }

  return (
    <>
      <header className="topbar">
        <div className="topbar-inner">
          <span className="brand">HAUP</span>
          <span className="section-name">ทะเบียนรถบริษัท</span>
        </div>
      </header>
      <main className="workspace">
        <div className="page-heading">
          <div>
            <h1>ข้อมูลรถยนต์</h1>
            <p>จัดการทะเบียน ยี่ห้อ รุ่น และหมายเหตุของรถบริษัท</p>
          </div>
          <button
            ref={addButton}
            className="button primary"
            disabled={load.pending || !!load.error}
            onClick={() => setAction({ kind: "add" })}
          >
            <span aria-hidden="true">＋</span> เพิ่มรถยนต์
          </button>
        </div>
        <p className="notice" role="status">
          {notice}
        </p>

        <section className="panel" aria-labelledby="car-list">
          <div className="toolbar">
            <div className="list-heading">
              <h2 id="car-list">รายการรถยนต์</h2>
              {!load.pending && !load.error && (
                <span className="count">
                  {needle ? `${visible.length} / ${cars.length}` : cars.length}{" "}
                  คัน
                </span>
              )}
            </div>
            <input
              aria-label="ค้นหารถยนต์"
              className="search"
              type="search"
              placeholder="ค้นหาทะเบียน ยี่ห้อ หรือรุ่น"
              value={query}
              onChange={(event) => setQuery(event.target.value)}
            />
          </div>

          {load.pending ? (
            <p className="empty" role="status">
              กำลังโหลดข้อมูลรถยนต์…
            </p>
          ) : load.error ? (
            <div className="empty">
              <p role="alert">{load.error}</p>
              <button
                className="button secondary"
                onClick={() => {
                  setLoad({ pending: true, error: "" });
                  setAttempt((value) => value + 1);
                }}
              >
                ลองอีกครั้ง
              </button>
            </div>
          ) : visible.length === 0 ? (
            <div className="empty">
              <strong>
                {needle ? "ไม่พบรถยนต์ที่ตรงกับคำค้น" : "ยังไม่มีข้อมูลรถยนต์"}
              </strong>
              <p>
                {needle
                  ? "ลองเปลี่ยนคำค้นหรือล้างช่องค้นหา"
                  : "เลือก “เพิ่มรถยนต์” เพื่อบันทึกรถคันแรก"}
              </p>
            </div>
          ) : (
            <div
              className="table-wrap"
              tabIndex={0}
              role="region"
              aria-label="ตารางข้อมูลรถยนต์"
            >
              <table>
                <thead>
                  <tr>
                    <th scope="col">ทะเบียนรถ</th>
                    <th scope="col">ยี่ห้อ / รุ่น</th>
                    <th scope="col">หมายเหตุ</th>
                    <th scope="col">อัปเดตล่าสุด</th>
                    <th scope="col">จัดการ</th>
                  </tr>
                </thead>
                <tbody>
                  {visible.map((car) => (
                    <tr key={car.id}>
                      <th scope="row">
                        <span className="plate">{car.registration}</span>
                      </th>
                      <td>
                        <strong>{car.brand}</strong>
                        <span className="model">{car.model}</span>
                      </td>
                      <td className="notes">{car.notes || "—"}</td>
                      <td className="date">
                        <time dateTime={car.updated_at}>
                          {dateFormat.format(new Date(car.updated_at))}
                        </time>
                      </td>
                      <td>
                        <div className="row-actions">
                          <button
                            className="button secondary small"
                            aria-label={`แก้ไขรถ ${car.registration}`}
                            onClick={() => setAction({ kind: "edit", car })}
                          >
                            แก้ไข
                          </button>
                          <button
                            className="button text-danger small"
                            aria-label={`ลบรถ ${car.registration}`}
                            onClick={() => setAction({ kind: "delete", car })}
                          >
                            ลบ
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      </main>
      {action && (
        <CarDialog
          action={action}
          onClose={() => setAction(null)}
          onSave={saved}
          onRemove={removed}
        />
      )}
    </>
  );
}
