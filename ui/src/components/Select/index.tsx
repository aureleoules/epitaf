import React from 'react';
import styles from './select.module.scss';

type Props = {
    title?: string
    children: any
    value?: string
    onChange?: any
    disabled?: boolean
}
export default function(props: Props) {
    return (
        <div className={styles.select}>
            <p>{props.title}</p>
            <select disabled={props.disabled} value={props.value} onChange={props.onChange}>
                {props.children}
            </select>
        </div>
    )
}