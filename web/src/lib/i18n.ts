import i18n from "i18next";
import { initReactI18next } from "react-i18next";

const resources = {
  en: {
    translation: {
      profile: {
        title: "Profile",
        subtitle: "Manage your account settings and preferences",
        information: {
          title: "Profile Information",
          subtitle: "Update your personal information and account details",
          fullName: "Full Name",
          username: "Username",
          email: "Email Address",
          fullNamePlaceholder: "Your full name",
          usernamePlaceholder: "Your username",
          emailPlaceholder: "your.email@example.com",
        },
        preferences: {
          title: "Preferences",
          language: "Language",
          theme: "Theme",
          timezone: "Timezone",
          languagePlaceholder: "English",
          timezonePlaceholder: "UTC",
        },
        accountStatus: {
          title: "Account Status",
          subtitle: "Current account information",
          accountType: "Account Type",
          status: "Account Status",
          lastUpdated: "Last Updated",
          lastSeen: "Last Seen",
          administrator: "Administrator",
          regularUser: "Regular User",
          active: "Active",
          disabled: "Disabled",
        },
        actions: {
          saveChanges: "Save Changes",
          cancel: "Cancel",
        },
        joined: "Joined",
        loading: "Loading profile...",
      },
      common: {
        english: "English",
        turkish: "Turkish",
      },
    },
  },
  tr: {
    translation: {
      profile: {
        title: "Profil",
        subtitle: "Hesap ayarlarınızı ve tercihlerinizi yönetin",
        information: {
          title: "Profil Bilgileri",
          subtitle: "Kişisel bilgilerinizi ve hesap detaylarınızı güncelleyin",
          fullName: "Ad Soyad",
          username: "Kullanıcı Adı",
          email: "E-posta Adresi",
          fullNamePlaceholder: "Adınız ve soyadınız",
          usernamePlaceholder: "Kullanıcı adınız",
          emailPlaceholder: "eposta@ornek.com",
        },
        preferences: {
          title: "Tercihler",
          language: "Dil",
          theme: "Tema",
          timezone: "Saat Dilimi",
          languagePlaceholder: "Türkçe",
          timezonePlaceholder: "UTC",
        },
        accountStatus: {
          title: "Hesap Durumu",
          subtitle: "Mevcut hesap bilgileri",
          accountType: "Hesap Türü",
          status: "Hesap Durumu",
          lastUpdated: "Son Güncellenme",
          lastSeen: "Son Görülme",
          administrator: "Yönetici",
          regularUser: "Normal Kullanıcı",
          active: "Aktif",
          disabled: "Devre Dışı",
        },
        actions: {
          saveChanges: "Değişiklikleri Kaydet",
          cancel: "İptal",
        },
        joined: "Katılım",
        loading: "Profil yükleniyor...",
      },
      common: {
        english: "English",
        turkish: "Türkçe",
      },
    },
  },
};

i18n.use(initReactI18next).init({
  resources,
  lng: "en",
  fallbackLng: "en",
  interpolation: {
    escapeValue: false,
  },
});

export default i18n;
