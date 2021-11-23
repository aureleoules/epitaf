import React, { useEffect, useState } from 'react';
import Button from '../../components/Button';
import styles from './signin.module.scss';

import history from '../../history';
import {ReactComponent as MicrosoftIcon} from '../../assets/svg/microsoft.svg';
import {ReactComponent as GitHubIcon} from '../../assets/svg/github.svg';
import {ReactComponent as Heart} from '../../assets/svg/heart.svg';
import {ReactComponent as AureleoulesLogo} from '../../assets/svg/aureleoules.svg';
import Client from '../../services/client';
import { useTranslation } from 'react-i18next';

import { ReactComponent as Loading } from '../../assets/svg/loading.svg';

export default function(props: any) {
    const {t} = useTranslation();

    const [authenticating, setAuthenticating] = useState<boolean>(false);

    function authenticate() {
        setAuthenticating(true);
        Client.Users.authenticateUrl().then(url => {
            window.location.replace(url);
        });
    }

    useEffect(() => {
        if (window.location.search.includes('token')) {
            const token = window.location.search.split('token=')[1];
            localStorage.setItem('jwt', token);
            window.location.replace('/');
        }
    }, []);

    return (
        <div className={["route " + styles.signin].join(" ")}>
            <div className={styles.container}>
                <h1>EPITAF</h1>
                <div>
                    <Button 
                        className={[authenticating ? styles.authenticating : "", styles.authenticate].join(" ")}
                        onClick={authenticate}
                        icon={authenticating ? Loading : MicrosoftIcon} 
                        large  
                        color="green" 
                        title={!authenticating ? "Sign in with Microsoft" : ""}
                    />
                </div>
                <div className={styles.infos}>
                    <div className={styles.links}>
                        <a rel="noopener noreferrer" target="_blank" href="https://github.com/aureleoules/epitaf">
                            <GitHubIcon/>
                        </a>
                        <a rel="noopener noreferrer" className={styles.aureleoules} target="_blank" href="https://aureleoules.com">
                            <AureleoulesLogo style={{width: 32}} />
                        </a>
                    </div>
                    <p>{t('Made with')} <Heart/> {t('by')} Aurèle</p>
                </div>
            </div>
        </div>
    )
}