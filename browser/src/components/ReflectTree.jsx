import React from "react";
import { Tree, Input } from 'antd';
import { connect } from "dva";
const TreeNode = Tree.TreeNode;
const Search = Input.Search;

/**/
const ReflectTree = (props) => {
    const buildTreeNode = (value, index) => {
        //console.log("buildTreeNode", value.key)
        let colorBegin = value.name.toLowerCase().indexOf(props.searchText);
        let beforeStr = value.name.substr(0, colorBegin);
        let hitStr = value.name.substr(colorBegin, props.searchText.length);
        let afterStr = value.name.substr(colorBegin + props.searchText.length);
        let title;
        if (colorBegin > -1) {
            title = <span style={{ fontSize: 16, fontFamily: "Consolas" }}>
                        {beforeStr}
                        <span style={{ color: "#CF0000"}}>{hitStr}</span>
                        {afterStr}
                    </span>;
        } else {
            title = <span style={{ fontSize: 16, fontFamily: "Consolas" }}>{value.name}</span>;
        }
        if (value.children) {
            return (
                <TreeNode key={value.key} title={title}>
                    {value.children.map(buildTreeNode)}
                </TreeNode>
            );
        }
        return <TreeNode key={value.key} title={title}/>;
    };

    const loadData = (e) => {
        props.dispatch({
            type: "reflectTree/changeDataPath",
            payload: {
                data: e.target.value
            }
        });
    }

    return (
        <div>
            <Input placeholder="host:port" onPressEnter={loadData} defaultValue={props.dataPath}
                style={{ width: "300px", marginTop: "20px", marginLeft: "5px", display: "block" }}/>
            <Search style={{ width: 300, marginTop: 10, marginLeft: 5, display: props.displayContent }} 
            placeholder="search" 
            onChange={ (event) => {props.dispatch({type:`reflectTree/search`, payload:{data:event}});} }/>
            <Tree style={{ display: props.displayContent }}
            expandedKeys={props.expandedKeys}
                    autoExpandParent={props.autoExpandParent}
                    onExpand={ (expandedKeys, {expanded, node}) => {props.dispatch({
                        type:`reflectTree/expand`, 
                        payload:{expandedKeys: expandedKeys, 
                            expanded: expanded,
                            node: node}
                        });} }>
                {props.treeNodeData.map(buildTreeNode)}
            </Tree>
        </div>
    );
};

ReflectTree.propTypes = {

};

function mapStateToProps({ reflectTree }) {
    return { ...reflectTree };
}

export default connect(mapStateToProps)(ReflectTree);