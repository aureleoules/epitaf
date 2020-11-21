import React from 'react';
import { useTranslation } from 'react-i18next';
import styles from './help.module.scss';

export default function(props: any) {

    const {t} = useTranslation();

    return (
        <div className={"route " + styles.help}>
            <h1>{t('Help')}</h1>
            <div>
                <p>{t("help_text")}</p>
                <ul>
                    <li>
                        {t('Create a')} <a target="_blank" rel="noopener noreferrer" href="https://github.com/aureleoules/epitaf/issues/new">GitHub issue</a>
                    </li>
                    <li>
                        {t('Contact me on Discord')}: <code>Nuf#8229</code>
                    </li>
                    <li>
                        {t('Send me an email')}: <a href="mailto:bug@epita.fr">bug@epita.fr</a>
                    </li>
                </ul>
            </div>
        </div>
    )
}