import { useState } from "react";

export default function RGBColorPicker({ onChange }) {
  const [color, setColor] = useState({ r: 0, g: 0, b: 0 });

  const handleChange = (channel, value) => {
    const newColor = { ...color, [channel]: parseInt(value, 10) };
    setColor(newColor);
    onChange?.(newColor);
  };

  return (
    <div style={{ padding: "1rem", background: "#222", color: "#fff" }}>
      {["r", "g", "b"].map((channel) => (
        <div key={channel}>
          <label>{channel.toUpperCase()}: {color[channel]}</label>
          <input
            type="range"
            min="0"
            max="255"
            value={color[channel]}
            onChange={(e) => handleChange(channel, e.target.value)}
          />
        </div>
      ))}
      <div style={{
        marginTop: "10px",
        width: "100px",
        height: "50px",
        backgroundColor: `rgb(${color.r}, ${color.g}, ${color.b})`,
        border: "1px solid #fff"
      }} />
    </div>
  );
}
