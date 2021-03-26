import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

import LanguageDetector from 'i18next-browser-languagedetector';
import moment from 'moment';
import fr from './translations/fr_FR.json';
import en from './translations/en_US.json';

import 'moment/locale/fr';

const resources = {
	en: {
		translation: en,
	},
	fr: {
		translation: fr,
	},
};

const languageDetector = new LanguageDetector(null, {
	lookupLocalStorage: 'language',
});

if (localStorage.getItem('language')?.startsWith('fr')) {
	moment.locale('fr');
} else {
	moment.locale('en');
}

i18n.use(initReactI18next) // passes i18n down to react-i18next
	.use(languageDetector)
	.init({
		resources,
		keySeparator: false, // we do not use keys in form messages.welcome

		interpolation: {
			escapeValue: false, // react already safes from xss
		},
	});

export default i18n;
