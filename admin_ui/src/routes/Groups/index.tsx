import React, { useEffect, useState } from 'react';

import GroupTree from '../../views/GroupTree';
import { useTranslation } from 'react-i18next';

export default function Groups(props: any) {

    const { t } = useTranslation();

    return (
        <>
            <h1>Groups</h1>
            <GroupTree/>
        </>
    )
}