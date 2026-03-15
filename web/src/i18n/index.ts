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

const LANGUAGE_STORAGE_KEY = "language";

const savedLanguage = localStorage.getItem(LANGUAGE_STORAGE_KEY);

// eslint-disable-next-line import-x/no-named-as-default-member
i18next.use(initReactI18next).init({
  resources,
  lng: savedLanguage && savedLanguage in resources ? savedLanguage : "en-US",
  fallbackLng: "en-US",
  interpolation: {
    escapeValue: false,
  },
});

i18next.on("languageChanged", (lng) => {
  localStorage.setItem(LANGUAGE_STORAGE_KEY, lng);
  document.documentElement.setAttribute("lang", lng);
});

// Set initial lang attribute
document.documentElement.setAttribute("lang", i18next.language || "en-US");

export default i18next;
