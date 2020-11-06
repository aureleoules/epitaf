import React, { ReactElement } from 'react';
import styles from './modal.module.scss';
import {ReactComponent as CrossIcon} from '../../assets/svg/cross.svg';

type Props = {
    children: ReactElement
    close?: any
}
export default function(props: Props) {
    return (
        <div className={styles.task}>
            <div onClick={props.close} className={styles.overlay}/>
            <div className={styles.container}>
                <CrossIcon onClick={props.close} className={styles.close}/>

                {props.children}
            </div>
        </div>
    )
}