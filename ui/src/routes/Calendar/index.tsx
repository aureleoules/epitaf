import React, { useEffect } from 'react';
import Calendar from '../../views/Calendar';
import styles from './calendar.module.scss';

export default function(props: any) {

    useEffect(() => {

    }, []);

    return (
        <div className={"route " + styles.calendar}>
            <Calendar />
        </div>
    )
}