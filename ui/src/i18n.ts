import i18n from "i18next";
import { initReactI18next } from "react-i18next";

import fr from './translations/fr_FR.json';
import en from './translations/en_US.json';

import LanguageDetector from 'i18next-browser-languagedetector';


// TODO : remove dayjs from project
import dayjs from "dayjs";
import moment from 'moment'
import 'moment/locale/fr'

import frLocale from 'dayjs/locale/fr';
import enLocale from 'dayjs/locale/en';

// the translations
// (tip move them in a JSON file and import them)
const resources = {
	en: {
		translation: en
	},
	fr: {
		translation: fr
	}
};

const languageDetector = new LanguageDetector(null, {
    lookupLocalStorage: 'language',
});

if(localStorage.getItem("language")?.startsWith("fr")) {
	dayjs.locale(frLocale);
	moment.locale('fr')
} else {
	dayjs.locale(enLocale);
	moment.locale('en');
}

i18n
	.use(initReactI18next) // passes i18n down to react-i18next
	.use(languageDetector)
	.init({
		resources,
		keySeparator: false, // we do not use keys in form messages.welcome

		interpolation: {
			escapeValue: false // react already safes from xss
		}
});

export default i18n;