import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { TreeNodeDatum } from 'react-d3-tree/lib/types/common';
import { Group } from '../../types/group';
import Client from '../../services/client';
import { convertGroupTreeToD3Tree } from '../../utils';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import Tree from 'react-d3-tree';
import history from '../../history';

export default function GroupTree(props: any) {
    const [groups, setGroups] = useState<Group | null>(null);
    const [node, setNode] = useState<TreeNodeDatum | null>(null);
    const { t } = useTranslation();

    useEffect(() => {
        Client.Groups.tree().then(groups => {
            setGroups(groups);
        }).catch(err => {
            if (err) throw err;
        });
    }, []);

    const [mousePos, setMousePos] = useState<any>({ x: 0, y: 0 });

    const foreignObjectProps = { width: 200, height: 50, x: 30, y: -20 };

    const renderForeignObjectNode = (props: any) => (
        <g onClick={() => setNode(props.nodeDatum)}>
            <circle style={{fill: "#505fb0"}} r={15}></circle>
            <foreignObject {...props.foreignObjectProps}>
                <h3 style={{margin: 0}}>{props.nodeDatum.name}</h3>
                <p style={{margin: 0}}>{t('slug')}: {props.nodeDatum.attributes.slug}</p>
            </foreignObject>
        </g>
    );
    
    return (
        <div onMouseOver={(e: any) => {
            if (node || e.target.nodeName !== "circle") return;
            setMousePos({
                x: e.pageX - document.body.scrollLeft,
                y: e.pageY - document.body.scrollTop,
            })
        }} style={{ width: '100%', height: 800 }}>
            {groups && <Tree
                translate={{
                    x: 300,
                    y: 50
                }}
                scaleExtent={{
                    min: 0.5,
                    max: 2
                }}
                separation={{
                    nonSiblings: 2,
                    siblings: 1
                }}
                collapsible={false}
                orientation={"vertical"}
                renderCustomNodeElement={(rd3tProps) =>
                    renderForeignObjectNode({ ...rd3tProps, foreignObjectProps })
                }
                data={convertGroupTreeToD3Tree(groups!)}
            />}
            {node && <Menu
                id="simple-menu"
                anchorReference="anchorPosition"
                anchorPosition={{
                    left: mousePos.x,
                    top: mousePos.y,
                }}
                anchorOrigin={{
                    vertical: 'top',
                    horizontal: 'left',
                }}
                transformOrigin={{
                    vertical: 'top',
                    horizontal: 'left',
                }}
                open={node !== null}
                onClose={() => setNode(null)}>
                <MenuItem onClick={() => history.push('/groups/' + node!.attributes!.uuid!)}>{t('Open')}</MenuItem>
                <MenuItem>{t('Move')}</MenuItem>
                <MenuItem>{t('Delete')}</MenuItem>
            </Menu>}
        </div>
    )
}