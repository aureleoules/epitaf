import React, {useEffect} from 'react';

import Client from '../../services/client';
import { getQueryVariable } from '../../utils';

import styles from './callback.module.scss';

import { ReactComponent as Loading } from '../../assets/svg/loading.svg';

export default function(props: any) {
    
    useEffect(() => {
        const code = getQueryVariable("code");
        if(!code) return;
        
        Client.Users.authenticate(code!);
    }, []) // es-lint-disable-line

    return (
        <div className={"route " + styles.callback}>
            <div className="loader loader--style2" title="1">
                <Loading/>
            </div>
        </div>
    )
}
