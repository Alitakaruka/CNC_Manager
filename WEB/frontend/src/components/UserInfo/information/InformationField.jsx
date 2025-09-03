import { useEffect, useState, createContext, useContext } from "react";
import "./InformationField.css";
import { setInformationHandler } from "../LogService.js"; // 👈 связываем глобально

export default function Information({ message, onClose, duration = 4000 }) {
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
      <button onClick={onClose}></button>
    </div>
  );
}

const InformationContext = createContext((msg) => {});

export function useInformation() {
  return useContext(InformationContext);
}

export function InformationProvider({ children }) {
  const [message, setMessage] = useState("");

  const showInformation = (msg) => {
    setMessage(msg);
    setTimeout(() => setMessage(""), 3000);
  };

  useEffect(() => {
    setInformationHandler(showInformation); // 👈 глобальная привязка
  }, []);

  return (
    <InformationContext.Provider value={showInformation}>
      {children}
      {message && <Information message={message} onClose={() => setMessage("")} />}
    </InformationContext.Provider>
  );
}
