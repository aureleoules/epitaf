import React, { useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { toast, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { Button, ControlLabel, Form, FormControl, FormGroup, HelpBlock, Input } from 'rsuite';
import { ReactComponent as Logo } from '../../assets/svg/bitcoin.svg';
import Client from '../../services/client';
import styles from './login.module.scss';

export default function Setup() {
	const { t } = useTranslation();

	const [email, setEmail] = useState<string>('');
	const [password, setPassword] = useState<string>('');

	function authenticateUser() {
		Client.Admins.authenticate(email, password)
			.then(() => {
				window.location.href = '/';
			})
			.catch((err) => {
				if (err) {
					toast(t('Incorrect username / password.'), {
						type: 'error',
					});
					throw err;
				}
			});
	}

	return (
		<div className={styles.login}>
			<h1>{t('Login')}</h1>
			<Form style={{marginTop: 25}} onSubmit={authenticateUser}>
				<Input onChange={v => setEmail(v)} style={{marginBottom: 15}} placeholder={t('Username')} name="name" type="email" />
				<Input onChange={v => setPassword(v)} style={{marginBottom: 15}} placeholder={t('Password')} name="password" type="password" />

				<Button type="submit" appearance="primary">{t('Login')}</Button>
			</Form>
		</div>
	);
}
