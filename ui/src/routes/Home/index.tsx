import React, { useEffect } from 'react';
import Tasks from '../../views/Tasks';
import Calendar from '../../views/Calendar';
import styles from './home.module.scss';


import ReactConfetti from 'react-confetti';
export default function(props: any) {

    useEffect(() => {
        setTimeout(() => {
            localStorage.removeItem("loginAnimation");
        }, 100);
    }, []);
    
    return <div className={"route " + styles.home}>
        {localStorage.getItem("loginAnimation") === "true" && <ReactConfetti numberOfPieces={300} recycle={false} />}
        
        <div className={styles.calendar}>
            <Calendar/>
        </div>
        
        <div className={styles.tasks}>
            <Tasks/>
        </div>
        
    </div>
}