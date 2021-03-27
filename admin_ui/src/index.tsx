import React from 'react';
import ReactDOM from 'react-dom';
import reportWebVitals from './reportWebVitals';
import Router from './Router';

import './i18n';
import './styles/styles.scss';
import 'animate.css';
import 'react-toastify/dist/ReactToastify.css';
import 'rsuite-table/dist/css/rsuite-table.css';
import 'rsuite/lib/styles/themes/dark/index.less';

ReactDOM.render(
	<React.StrictMode>
		<Router />
	</React.StrictMode>,
	document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
