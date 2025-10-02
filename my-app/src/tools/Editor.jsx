import { useQuill } from "react-quilljs";
import "quill/dist/quill.snow.css";

export default function Editor({ value, onChange }) {
  const { quill, quillRef } = useQuill({
    theme: "snow",
    modules: {
      toolbar: [
        ["bold", "italic", "underline", "strike"],
        [{ header: 1 }, { header: 2 }],
        [{ list: "ordered" }, { list: "bullet" }],
        ["link", "image"],
        ["clean"],
      ],
    },
  });

  // sync external value → editor
  if (quill && value && quill.root.innerHTML !== value) {
    quill.root.innerHTML = value;
  }

  // sync editor → parent
  if (quill) {
    quill.on("text-change", () => {
      onChange(quill.root.innerHTML);
    });
  }

  return <div ref={quillRef} style={{ height: "40vh" }} />;
}
