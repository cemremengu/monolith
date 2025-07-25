import i18next from "i18next";
import { initReactI18next } from "react-i18next";

import enTranslation from "./en/translation.json";
import trTranslation from "./tr/translation.json";

const resources = {
  en: {
    translation: enTranslation,
  },
  tr: {
    translation: trTranslation,
  },
};

// eslint-disable-next-line import-x/no-named-as-default-member
i18next.use(initReactI18next).init({
  resources,
  lng: "en",
  fallbackLng: "en",
  interpolation: {
    escapeValue: false,
  },
});

export default i18next;
