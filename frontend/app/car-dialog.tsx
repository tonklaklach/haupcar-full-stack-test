import {
  useEffect,
  useRef,
  useState,
  type FormEvent,
  type ChangeEvent,
} from "react";
import { carsAPI, type Car, type CarInput } from "../lib/cars";
import { formatRegistration, registrationError } from "../lib/registration";

export type DialogAction =
  | { kind: "add" }
  | { kind: "edit"; car: Car }
  | { kind: "delete"; car: Car };

type Props = {
  action: DialogAction;
  onClose: () => void;
  onSave: (car: Car) => void;
  onRemove: (id: number) => void;
};

export function CarDialog({ action, onClose, onSave, onRemove }: Props) {
  const dialog = useRef<HTMLDialogElement>(null);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const [plateError, setPlateError] = useState("");
  const car = action.kind === "add" ? undefined : action.car;
  const deleting = action.kind === "delete";
  const title = deleting
    ? "ลบข้อมูลรถยนต์"
    : car
      ? "แก้ไขข้อมูลรถยนต์"
      : "เพิ่มรถยนต์";

  useEffect(() => {
    const element = dialog.current;
    element?.showModal();
    return () => element?.close();
  }, []);

  function updateRegistration(field: HTMLInputElement) {
    const value = field.value;
    const formatted = formatRegistration(value);
    if (formatted && formatted !== value) {
      const caret = field.selectionStart ?? value.length;
      const beforeCaret = value.slice(0, caret);
      const nextCaret =
        formatRegistration(beforeCaret)?.length ??
        beforeCaret.replace(/\p{White_Space}/gu, "").length;
      field.value = formatted;
      field.setSelectionRange(nextCaret, nextCaret);
    }
    setPlateError("");
  }

  function onRegistrationChange(event: ChangeEvent<HTMLInputElement>) {
    if ((event.nativeEvent as InputEvent).isComposing) return;
    updateRegistration(event.currentTarget);
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (busy) return;
    const form = new FormData(event.currentTarget);
    setError("");
    setBusy(true);
    try {
      if (action.kind === "delete") {
        await carsAPI.remove(action.car.id);
        onRemove(action.car.id);
      } else {
        const input: CarInput = {
          registration: String(form.get("registration") ?? "").trim(),
          brand: String(form.get("brand") ?? "").trim(),
          model: String(form.get("model") ?? "").trim(),
          notes: String(form.get("notes") ?? "").trim(),
        };
        if (!input.registration || !input.brand || !input.model) {
          setError("กรุณาระบุทะเบียนรถ ยี่ห้อ และรุ่น");
          return;
        }
        const formatted = formatRegistration(input.registration);
        if (!formatted) {
          setPlateError(registrationError);
          const field = event.currentTarget.elements.namedItem("registration");
          if (field instanceof HTMLInputElement) field.focus();
          return;
        }
        onSave(
          await carsAPI.save({ ...input, registration: formatted }, car?.id),
        );
      }
      onClose();
    } catch (reason) {
      setError(
        reason instanceof Error ? reason.message : "บันทึกข้อมูลไม่สำเร็จ",
      );
    } finally {
      setBusy(false);
    }
  }

  return (
    <dialog
      ref={dialog}
      aria-labelledby="dialog-title"
      onCancel={(event) => {
        event.preventDefault();
        if (!busy) onClose();
      }}
    >
      <form onSubmit={submit}>
        <h2 id="dialog-title">{title}</h2>
        {deleting ? (
          <p className="delete-description">
            ต้องการลบรถทะเบียน <strong>{car?.registration}</strong> หรือไม่?
            ข้อมูลที่ลบแล้วจะไม่สามารถกู้คืนได้
          </p>
        ) : (
          <fieldset disabled={busy}>
            <label>
              ทะเบียนรถ <span aria-hidden="true">*</span>
              <input
                autoFocus
                name="registration"
                placeholder="เช่น 4ขห-798 หรือ ขก-1234"
                defaultValue={
                  car
                    ? (formatRegistration(car.registration) ?? car.registration)
                    : ""
                }
                onChange={onRegistrationChange}
                onCompositionEnd={(event) =>
                  updateRegistration(event.currentTarget)
                }
                onBlur={(event) => {
                  const value = event.currentTarget.value.trim();
                  setPlateError(
                    value && !formatRegistration(value)
                      ? registrationError
                      : "",
                  );
                }}
                aria-invalid={plateError ? true : undefined}
                aria-describedby={plateError ? "registration-error" : undefined}
                required
                maxLength={30}
              />
              {plateError && (
                <span
                  id="registration-error"
                  className="field-error"
                  role="alert"
                >
                  {plateError}
                </span>
              )}
            </label>
            <div className="field-row">
              <label>
                ยี่ห้อ <span aria-hidden="true">*</span>
                <input
                  name="brand"
                  placeholder="เช่น Toyota"
                  defaultValue={car?.brand}
                  required
                  maxLength={60}
                />
              </label>
              <label>
                รุ่น <span aria-hidden="true">*</span>
                <input
                  name="model"
                  placeholder="เช่น Yaris"
                  defaultValue={car?.model}
                  required
                  maxLength={80}
                />
              </label>
            </div>
            <label>
              หมายเหตุ
              <textarea
                name="notes"
                rows={3}
                placeholder="ข้อมูลเพิ่มเติม (ถ้ามี)"
                defaultValue={car?.notes}
                maxLength={1000}
              />
            </label>
          </fieldset>
        )}
        {error && (
          <p className="form-error" role="alert">
            {error}
          </p>
        )}
        <div className="dialog-actions">
          <button
            autoFocus={deleting}
            className="button secondary"
            type="button"
            disabled={busy}
            onClick={onClose}
          >
            ยกเลิก
          </button>
          <button
            className={`button ${deleting ? "danger" : "primary"}`}
            disabled={busy}
          >
            {busy
              ? "กำลังดำเนินการ…"
              : deleting
                ? "ยืนยันการลบ"
                : "บันทึกข้อมูล"}
          </button>
        </div>
      </form>
    </dialog>
  );
}
