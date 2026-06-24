import * as ToastPrimitive from "@radix-ui/react-toast";
import {
  createContext,
  type ReactNode,
  useCallback,
  useContext,
  useMemo,
  useRef,
  useState,
} from "react";
import { Toast, type ToastIntent } from "~/components/atoms/toast/Toast";

const DEFAULT_DURATION = 4000;

interface ToastItem {
  id: number;
  type: ToastIntent;
  title: string;
  message: string;
}

interface ToastContextValue {
  success: (message?: string, title?: string) => void;
  error: (message?: string, title?: string) => void;
}

const ToastContext = createContext<ToastContextValue | null>(null);

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastItem[]>([]);
  const counterRef = useRef(0);

  const remove = useCallback((id: number) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const add = useCallback(
    (type: ToastIntent, title: string, message: string) => {
      const id = ++counterRef.current;
      setToasts((prev) => [...prev, { id, type, title, message }]);
    },
    [],
  );

  const value = useMemo<ToastContextValue>(
    () => ({
      success: (message = "", title = "Success") =>
        add("success", title, message),
      error: (message = "", title = "Oops something went wrong") =>
        add("error", title, message),
    }),
    [add],
  );

  return (
    <ToastContext.Provider value={value}>
      <ToastPrimitive.Provider>
        {children}
        {toasts.map((toast) => (
          <Toast
            key={toast.id}
            intent={toast.type}
            title={toast.title}
            description={toast.message}
            duration={DEFAULT_DURATION}
            onClose={() => remove(toast.id)}
          />
        ))}
        <ToastPrimitive.Viewport className="fixed right-4 bottom-4 z-50 flex w-80 max-w-full flex-col gap-2 outline-none" />
      </ToastPrimitive.Provider>
    </ToastContext.Provider>
  );
}

export function useToast() {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error("useToast must be used within a ToastProvider");
  }
  return context;
}
