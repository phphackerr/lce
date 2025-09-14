import tippy from "tippy.js";
import "tippy.js/dist/tippy.css";
// import "tippy.js/themes/light-border.css";
import "tippy.js/themes/material.css";
import "tippy.js/animations/scale.css";

export const defaultOptions = {
  placement: "auto",
  animation: "scale",
  arrow: true,
  delay: [200, 100],
  allowHTML: true,
  // theme: "light-border",
  theme: "material",
  inertia: true,
  popperOptions: {
    modifiers: [
      {
        name: "preventOverflow",
        options: {
          padding: 8,
        },
      },
    ],
  },
};

// @ts-ignore
export function tt(node, params) {
  const options = { ...defaultOptions, ...params };
  const instance = tippy(node, options);
  return {
    // @ts-ignore
    update(newParams) {
      // @ts-ignore
      instance.setProps({ ...defaultOptions, ...newParams });
    },
    destroy() {
      // @ts-ignore
      instance.destroy();
    },
  };
}

if (typeof window !== "undefined") {
  const style = document.createElement("style");
  style.textContent = `
    .tippy-box {
      user-select: text !important;
      z-index: 100000000 !important;
    }
  `;
  document.head.appendChild(style);
}
