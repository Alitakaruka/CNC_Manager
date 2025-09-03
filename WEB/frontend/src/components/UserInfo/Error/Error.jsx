import { useEffect, useState, createContext, useContext } from "react";
import "./Error.css";
import { setErrorHandler } from "../LogService"; // 👈 связываем глобально

export default function Error({ message, onClose, duration = 4000 }) {
  useEffect(() => {
    const timer = setTimeout(() => {
      onClose();
    }, duration);
    return () => clearTimeout(timer);
  }, [onClose, duration]);

  if (!message) return null;

  return (
    <div className="toast">
      <span>{message}</span>
      <button onClick={onClose}>×</button>
    </div>
  );
}

const ErrorContext = createContext((msg) => {});

export function useError() {
  return useContext(ErrorContext);
}

export function ErrorProvider({ children }) {
  const [message, setMessage] = useState("");

  const showError = (msg) => {
    setMessage(msg);
    setTimeout(() => setMessage(""), 3000);
  };

  useEffect(() => {
    setErrorHandler(showError); // 👈 глобальная привязка
  }, []);

  return (
    <ErrorContext.Provider value={showError}>
      {children}
      {message && <Error message={message} onClose={() => setMessage("")} />}
    </ErrorContext.Provider>
  );
}
