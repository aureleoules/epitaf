import dayjs from 'dayjs';
import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { RotateSpinner } from "react-spinners-kit";
import chronosMapping from '../../assets/data/chronos_mapping.json';
import Client from '../../services/client';
import { Calendar } from '../../types/calendar';
import { IDictionary } from '../../types/dictionnary';
import styles from './calendar.module.scss';


const colors: IDictionary<string> = chronosMapping;

type Props = {
}
export default function(props: Props) {
    const {t} = useTranslation();
    const [calendar, setCalendar] = useState<Array<Calendar> | null>(null);
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

    function getDays() {
        const map = new Map();
        calendar?.forEach((item) => {
             const key = new Date(item.startDate).setHours(0, 0, 0, 0);
             const collection = map.get(key);
             if (!collection) {
                 map.set(key, [item]);
             } else {
                 collection.push(item);
             }
        });

        const o = Object.fromEntries(map); 

        return Object.keys(o).sort().reduce(
            (obj, key) => { 
                (obj as any)[key] = o[key]; 
                return obj;
            }, 
            {}
        );
    }

    return (
        <div className={styles.calendar}>
            <h1>{t('Calendar')}</h1>
            
            {(calendar?.length === 0 && fetched) && <p>{t('Unavailable')}</p>}
            
            {!fetched && <div style={{position: "absolute", left: "40%", top: "45%"}}>
                <RotateSpinner size={50} color="var(--primary)"/>
            </div>}
            {calendar?.length !== 0 && fetched && <>
                <div className={[styles.classes].join(" ")}>
                    {Object.keys(getDays()).map((d: any, i: number) => {
                        const day = (getDays() as any)[d];
                        return (
                            <div key={i} className={styles.day}>
                                <h2>{dayjs(new Date(parseInt(d))).format("dddd DD MMMM")}</h2>
                                {day.sort((a: any, b: any) => new Date(a.startDate).getTime() - new Date(b.startDate).getTime() > 0 ? 1 : -1).map((c: any, i: number) => {
                                    const duration = (new Date(c.endDate).getTime() - new Date(c.startDate).getTime()) / 1000 / 60;

                                    return (
                                       <div key={i} className={styles.class}>
                                           <div className={styles.date}>
                                               <h3>{dayjs(c.startDate).format("HH:mm")}</h3>
                                               <p>{Math.floor(duration / 60) + "h" + ((duration % 60) > 0 ? (duration % 60) : "")}</p>
                                           </div>
                                           <span 
                                               style={{
                                                   backgroundColor: colors[c.name.toString()]
                                               }} 
                                               className={styles.separator}
                                           />
                                           <div className={styles.content}>
                                               <span>{(c.rooms && c.rooms.length > 0) ? c.rooms.map((x: any) => x.name).join(', ') : ""}</span>
                                               <p>{c.name ? c.name : ""}</p>
                                           </div>
                                       </div>
                                   )
                                })}
                            </div>
                        )
                    })}
                </div>
            </>}
        </div>
    )
}
