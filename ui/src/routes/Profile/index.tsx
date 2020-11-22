import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import Client from '../../services/client';
import { User } from '../../types/user';
import styles from './profile.module.scss';

export default function(props: any) {

    const {t} = useTranslation();
    const [profile, setProfile] = useState<User | null>(null);

    useEffect(() => {
        Client.Users.me().then(profile => {
            setProfile(profile);
        }).catch(err => {
            if(err) throw err;
        });
    }, []);

    return (
        <div className={"route " + styles.profile}>
            <h1>{t('Profile')}</h1>
            {profile && <div className={styles.container}>
                <h2>{profile.name}</h2>
                <p>
                    <span>
                        {t('Email')} :
                    </span>
                    {profile.email}</p>
                <p>
                    <span>
                        {t('Promotion')} :
                    </span>
                    {profile.promotion}</p>
                <p>
                    <span>
                        {t('Semester')} :
                    </span>
                    {profile.semester}</p>
                <p>
                    <span>
                        {t('Class')} :
                    </span>
                    {profile.class}</p>
                <p>
                    <span>
                        {t('Region')} :
                    </span>
                    {profile.region}</p>
            </div>}
        </div>
    )
}