import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { TreeNodeDatum } from 'react-d3-tree/lib/types/common';
import Tree from 'react-d3-tree';
import { Button, ControlLabel, Dropdown, Form, FormControl, FormGroup, Modal } from 'rsuite';
import { createRef } from 'reactn';
import { Group } from '../../types/group';
import Client from '../../services/client';
import { convertGroupTreeToD3Tree } from '../../utils';
import styles from './grouptree.module.scss';
import history from '../../history';

export default function GroupTree(props: any) {
	const [groups, setGroups] = useState<Group | null>(null);
	const [node, setNode] = useState<TreeNodeDatum | null>(null);
	const [group, setGroup] = useState<Group>({});
	const [createGroupModal, setCreateGroupModal] = useState<boolean>(false);
	const { t } = useTranslation();

	useEffect(() => {
		fetchTree();
	}, []);

	function fetchTree() {
		Client.Groups.tree().then(groups => {
			setGroups(groups);
		}).catch(err => {
			if (err) throw err;
		});
	}
    
	function createGroup() {
		Client.Groups.create(group.parent_id!, group).then(() => {
			fetchTree();
			setCreateGroupModal(false);
		}).catch(err => {
			if (err) throw err;
		});
	}

	function removeGroup(id: string) {
		Client.Groups.delete(id).then(() => {
			fetchTree();
		}).catch(err => {
			if (err) throw err;
		});
	}

	const foreignObjectProps = { width: 200, height: 200, x: 30, y: -20 };

	const renderForeignObjectNode = (props: any) => {
		const r = createRef<any>();
		return (
			<>
				<g onClick={() => setNode(props.nodeDatum)}>
					<circle onClick={e => r.current.handleClick(e)} style={{fill: '#505fb0'}} r={8} />
					<foreignObject {...props.foreignObjectProps}>
						<Dropdown
							ref={r}
							className={styles.dropdown}
							renderTitle={() => 
								<>
									<h4 style={{margin: 0}}>{props.nodeDatum.name}</h4>
									<p 
										style={{margin: 0, fontSize: 12, width: 350}}
									>
										{`${t('slug')}: ${props.nodeDatum.attributes.slug}`}
									</p>
								</>}
							title="edit"
						>
							<Dropdown.Item onSelect={() => history.push(`/groups/${props.nodeDatum.attributes.id}`)}>{t('Open')}</Dropdown.Item>
							<Dropdown.Item onSelect={
								() => {
									console.log(props.nodeDatum.attributes);
									setGroup(g => ({...g, parent_id: props.nodeDatum.attributes.id}));
									setCreateGroupModal(true);
								}
							}
							>{t('Add subgroup')}
							</Dropdown.Item>
							<Dropdown.Item onSelect={() => removeGroup(props.nodeDatum.attributes.id)}>{t('Remove')}</Dropdown.Item>
						</Dropdown>
					</foreignObject>
				</g>
			</>
		);
	};
	
	return (
		<div
			style={{ width: '100%', height: 800 }}
		>
			{groups && <Tree
				translate={{
					x: 300,
					y: 100
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
				orientation="vertical"
				renderCustomNodeElement={(rd3tProps) =>
					renderForeignObjectNode({ ...rd3tProps, foreignObjectProps })
				}
				data={convertGroupTreeToD3Tree(groups!)}
			/>}
			<Modal
				show={createGroupModal}
				close={() => setCreateGroupModal(false)}
				width={500}
				onHide={() => setCreateGroupModal(false)} 
			>
				<Modal.Header>
					<Modal.Title>{t('Create group')}</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					<Form fluid onSubmit={createGroup}>
						<FormGroup>
							<ControlLabel>{t('Name')}</ControlLabel>
							<FormControl
								autoFocus
								value={group.name}
								onChange={(v) => setGroup(g => ({...g, name: v}))}
								placeholder={t('Name')}
							/>
						</FormGroup>
					</Form>
				</Modal.Body>
				<Modal.Footer>
					<Button
						appearance="primary"
						type="submit"
						onClick={createGroup}
						disabled={!group.name}
					>
						{t('Create group')}
					</Button>
					<Button onClick={() => setCreateGroupModal(false)} appearance="subtle">{t('Cancel')}</Button>
				</Modal.Footer>
			</Modal>
		</div>
	);
}