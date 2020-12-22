import React, { useEffect } from 'react';
import { AddIcon, DataGrid } from '@material-ui/data-grid';
import Client from '../../services/client';
import { Fab, makeStyles } from '@material-ui/core';
import { useTranslation } from 'react-i18next';

const useStyles = makeStyles((theme) => ({
    root: {
        backgroundColor: theme.palette.background.paper,
        width: 500,
        position: 'relative',
        minHeight: 200,
    },
    fab: {
        position: 'absolute',
        bottom: theme.spacing(2),
        right: theme.spacing(2),
    }
}));

export default function Users(props: any) {
    
    const classes = useStyles();
    const {t} = useTranslation();
    
    const columns = [
        { field: 'id', headerName: 'ID', width: 70 },
        { field: 'firstName', headerName: 'First name', width: 130 },
        { field: 'lastName', headerName: 'Last name', width: 130 },
        {
            field: 'age',
            headerName: 'Age',
            type: 'number',
            width: 90,
        },
        {
            field: 'fullName',
            headerName: 'Full name',
            description: 'This column has a value getter and is not sortable.',
            sortable: false,
            width: 160,
            valueGetter: (params: any) =>
                `${params.getValue('firstName') || ''} ${params.getValue('lastName') || ''}`,
        },
    ];

    useEffect(() => {
        Client.Users.list().then(users => {

        }).catch(err => {
            if (err) throw err;
        });
    }, []);

    const rows = [
        { id: 1, lastName: 'Snow', firstName: 'Jon', age: 35 },
        { id: 2, lastName: 'Lannister', firstName: 'Cersei', age: 42 },
        { id: 3, lastName: 'Lannister', firstName: 'Jaime', age: 45 },
        { id: 4, lastName: 'Stark', firstName: 'Arya', age: 16 },
        { id: 5, lastName: 'Targaryen', firstName: 'Daenerys', age: null },
        { id: 6, lastName: 'Melisandre', firstName: null, age: 150 },
        { id: 7, lastName: 'Clifford', firstName: 'Ferrara', age: 44 },
        { id: 8, lastName: 'Frances', firstName: 'Rossini', age: 36 },
        { id: 19, lastName: 'Roxie', firstName: 'Harvey', age: 65 },
        { id: 911, lastName: 'Roxie', firstName: 'Harvey', age: 65 },
        { id: 9984, lastName: 'Roxie', firstName: 'Harvey', age: 65 },
        { id: 9897, lastName: 'Roxie', firstName: 'Harvey', age: 65 },
        { id: 994, lastName: 'Roxie', firstName: 'Harvey', age: 65 },
        { id: 9654, lastName: 'Roxie', firstName: 'Harvey', age: 65 },
        { id: 9321, lastName: 'Roxie', firstName: 'Harvey', age: 65 },
        { id: 987, lastName: 'Roxie', firstName: 'Harvey', age: 65 },
        { id: 97, lastName: 'Roxie', firstName: 'Harvey', age: 65 },
        { id: 99, lastName: 'Roxie', firstName: 'Harvey', age: 65 },
        { id: 859, lastName: 'Roxie', firstName: 'Harvey', age: 65 },
        { id: 99827, lastName: 'Roxie', firstName: 'Harvey', age: 65 },
        { id: 91111, lastName: 'Roxie', firstName: 'Harvey', age: 65 },
    ];


    return (
        <>
            <h1>Users</h1>
            <div style={{ height: 800, width: "100%" }}>
                <DataGrid rows={rows} columns={columns} pageSize={15} checkboxSelection />
            </div>

            <Fab aria-label={t('Add user')} className={classes.fab} color="primary">
                <AddIcon/>
            </Fab>
        </>
    );
}