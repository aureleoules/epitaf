import React from 'react';
import styles from './checkbox.module.scss';

type Props = {
    title?: string
    onChange?: any
    checked?: boolean
    disabled?: boolean
}
export default function Checkbox(props: Props) {
    return (
        <div className={[styles.checkbox, props.disabled ? styles.disabled : ""].join(" ")}>
            <label className={styles.container}>
                <p>
                    {props.title}
                </p>
                <input checked={props.checked} onChange={props.onChange} type="checkbox"/>
                <span className={styles.checkmark}></span>
            </label>
        </div>
    )
}