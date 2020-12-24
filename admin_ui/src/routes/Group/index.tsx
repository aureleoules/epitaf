import React, { useEffect } from 'react';
import Client from '../../services/client';

export default function Group(props: any) {
    console.log(props);
    useEffect(() => {
        Client.Groups.get(props.match.params.id).then(group => {
            console.log(group);
        }).catch(err => {
            if(err) throw err;
        });
    }, [props.match.params.id]);
    
    return (
        <>
            <h1>Group</h1>
        </>
    )
}