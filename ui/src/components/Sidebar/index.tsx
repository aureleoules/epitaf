import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

import {ReactComponent as TaskIcon} from '../../assets/svg/task.svg';
import {ReactComponent as HomeIcon} from '../../assets/svg/home.svg';
import {ReactComponent as CalendarIcon} from '../../assets/svg/calendar.svg';
import {ReactComponent as UserIcon} from '../../assets/svg/user.svg';
import {ReactComponent as MenuIcon} from '../../assets/svg/menu.svg';
import {ReactComponent as MoonIcon} from '../../assets/svg/moon.svg';
import {ReactComponent as SunIcon} from '../../assets/svg/sun.svg';
import {ReactComponent as FranceIcon} from '../../assets/svg/france.svg';
import {ReactComponent as USAIcon} from '../../assets/svg/usa.svg';

import styles from './sidebar.module.scss';
import history from '../../history';
import { useTranslation } from 'react-i18next';
import Client from '../../services/client';
import { User } from '../../types/user';
import Button from '../Button';
import { getTheme, setTheme, getLanguage, setLanguage } from '../../utils';
import i18n from '../../i18n';

function isCurrentRoute(current: string, route: string): boolean {
    return current.split("/")[1] === route.split("/")[1]
}

type Props = {}
export default function(props: Props) {

    const [user, setUser] = useState<User | null>(null);
    const [route, setRoute] = useState<string>(history.location.pathname);
    const [menu, setMenu] = useState<boolean>(false);
    
    const {t} = useTranslation();

    const routes = [
        {
            title: t('Home'),
            path: "/",
            icon: HomeIcon
        },
        {
            title: t('Tasks'),
            path: "/tasks",
            icon: TaskIcon
        },
        {
            title: t('Calendar'),
            path: "/calendar",
            icon: CalendarIcon
        }
    ]

    const bottom_routes = [
        {
            title: t('Profile'),
            path: "/me",
            icon: UserIcon
        }
    ]

    function signout() {
        localStorage.setItem("jwt", "");
        window.location.replace("/");
    }

    useEffect(() => {
        Client.Users.me().then(user => {
            setUser(user);
        }).catch(err => {
            if(err) throw err;
        });
        
        history.listen(r => {
            setRoute(r.pathname);
        });
    }, []);

    function toggleMenu() {
        setMenu(!menu);
    }

    function toggleTheme() {
        if(getTheme() === "dark") setTheme("light");
        else setTheme("dark");

        window.dispatchEvent(new Event("render"));
    }

    function toggleLanguage() {
        let lang = getLanguage() === "fr-FR" ? "en-US": "fr-FR";
        
        setLanguage(lang);
        i18n.changeLanguage(lang);
    }

    return (
        <div className={styles.sidebar}>
            <div className={[styles["container-content"], menu ? styles.isOpen : styles.isClosed].join(" ")}>
                <MenuIcon className={styles.burger} onClick={toggleMenu}/>
                <div className={[styles.container].join(" ")}>
                    <div className={styles.logo}>
                        <p>EPITAF</p>
                        {getTheme() === "light" && <button className={styles.c} onClick={toggleTheme}>
                            <SunIcon/>
                        </button>}
                        {getTheme() === "dark" && <button className={styles.c} onClick={toggleTheme}>
                            <MoonIcon/>
                        </button>}
                        <div className={styles.languages}>
                            {getLanguage() === "fr-FR" && <button onClick={toggleLanguage} className={styles.c}>
                                <FranceIcon/>
                            </button>}
                            {getLanguage() === "en-US" && <button onClick={toggleLanguage} className={styles.c}>
                                <USAIcon/>
                            </button>}
                        </div>
                        
                    </div>
                    <div className={styles.wrapper}>
                        <div className={styles.user}>
                            <div 
                                className={styles.avatar} 
                                style={{background: `url(${"https://photos.cri.epita.fr/square/" + user?.login})`}}
                            />
                            {user && <div className={styles.content}>
                                <h4>{user.name}</h4>
                                {!user.teacher && <>
                                    {!user.class ? <p>{t('sidebar_user_region', {semester: user.semester, region: user.region})}</p> : 
                                        <p>{t('sidebar_user_class', {semester: user.semester, class: user.class})}</p>}
                                </>}
                                {user.teacher && <p>{t('Teacher')}</p>}
                            </div>}
                        </div>

                        <div className={styles.routes}>
                            <ul>
                                {routes.map((r, i) => (
                                    <li key={i} className={isCurrentRoute(route, r.path) ? styles.active : ""}>
                                        <Link onClick={() => setMenu(false)} to={r.path}>
                                            <r.icon/>
                                            {r.title}
                                        </Link>
                                    </li>
                                ))}
                            </ul>
                            <ul className={styles.bottom}>
                                {bottom_routes.map((r, i) => (
                                    <li key={i} className={isCurrentRoute(route, r.path) ? styles.active : ""}>
                                        <Link onClick={() => setMenu(false)} to={r.path}>
                                            <r.icon/>
                                            {r.title}
                                        </Link>
                                    </li>
                                ))}
                                <div className={styles.signout}>
                                    <Button onClick={signout} color="red" title={t('Sign out')}/>
                                </div>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}