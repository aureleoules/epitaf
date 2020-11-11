import React, { useEffect, useState } from 'react';
import Client from '../../services/client';
import { Calendar } from '../../types/calendar';
import styles from './calendar.module.scss';
import dayjs from 'dayjs';
import chronosMapping from '../../assets/data/chronos_mapping.json';
import { IDictionary } from '../../types/dictionnary';
import { RotateSpinner   } from "react-spinners-kit";
import { useTranslation } from 'react-i18next';


const colors: IDictionary<string> = chronosMapping;

type Props = {
}
export default function(props: Props) {
    const {t} = useTranslation();
    const [calendar, setCalendar] = useState<Calendar | null>(null);
    const [fetched, setFetched] = useState<boolean>(false);

    useEffect(() => {
        Client.Users.calendar().then(cal => {
            setCalendar(cal);
            setFetched(true);
        }).catch(err => {
            if(err) {
                setFetched(true);
                throw err
            };
        });
    }, []);

    return (
        <div className={styles.calendar}>
            <h1>{t('Calendar')}</h1>
            
            {(!calendar?.days && fetched) && <p>{t('Unavailable')}</p>}
            
            {!fetched && <div style={{position: "absolute", left: "40%", top: "45%"}}>
                <RotateSpinner size={50} color="#572ce8"/>
            </div>}
            {calendar?.days && <>
                
                <div className={[styles.classes].join(" ")}>
                    {calendar?.days.map((d, i) => (
                        <div key={i} className={styles.day}>
                            <h2>{dayjs(d.date).format("DD MMMM")}</h2>

                            {d.courses.map((c, i) => (
                                <div key={i} className={styles.class}>
                                    <div className={styles.date}>
                                        <h3>{dayjs(c.start_date).format("HH:mm")}</h3>
                                        <p>{Math.floor(c.duration / 60) + "h" + ((c.duration % 60) > 0 ? (c.duration % 60) : "")}</p>
                                    </div>
                                    <span style={{
                                        backgroundColor: colors[c.name.toString()]
                                    }} className={styles.separator}/>
                                    <div className={styles.content}>
                                        <span>{c.rooms ? c.rooms[0] : ""}</span>
                                        <p>{c.name ? c.name : ""}</p>
                                    </div>
                                </div>
                            ))}
                        </div>
                    ))}
                </div>
            </>}
        </div>
    )
}