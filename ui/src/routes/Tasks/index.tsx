import React from 'react';
import Tasks from '../../views/Tasks';
import styles from './tasks.module.scss';

export default function(props: any) {

    return (
        <div className={styles.tasks + " route"}>
            <Tasks {...props}/>
        </div>
    )
}