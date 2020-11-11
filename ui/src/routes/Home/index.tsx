import React from 'react';
import Tasks from '../../views/Tasks';
import Calendar from '../../views/Calendar';
import styles from './home.module.scss';

export default function(props: any) {
    return <div className={"route " + styles.home}>
        <div className={styles.calendar}>
            <Calendar/>
        </div>
        
        <div className={styles.tasks}>
            <Tasks/>
        </div>
        
    </div>
}