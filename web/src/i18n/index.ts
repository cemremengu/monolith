import i18next from "i18next";
import { initReactI18next } from "react-i18next";

import enUSTranslation from "./en-US/translation.json";
import trTRTranslation from "./tr-TR/translation.json";

const resources = {
  "en-US": {
    translation: enUSTranslation,
  },
  "tr-TR": {
    translation: trTRTranslation,
  },
};

// eslint-disable-next-line import-x/no-named-as-default-member
i18next.use(initReactI18next).init({
  resources,
  lng: "en-US",
  fallbackLng: "en-US",
  interpolation: {
    escapeValue: false,
  },
});

// Ensure the HTML lang attribute is set to the full BCP 47 code
i18next.on("languageChanged", (lng) => {
  document.documentElement.setAttribute("lang", lng);
});

// Set initial lang attribute
document.documentElement.setAttribute("lang", i18next.language || "en-US");

export default i18next;
