import * as request from '../utils/request';

function getData(key) {
    const parseData = (res) => {
        let data = null;

        if ("err" in res) {
            console.log(res);
        } else {
            data = res.data;
        }
        return data;
    }

    return request.default(`${request.HTTPPREFIX}query?key=${key}`)
        .then(parseData)
        .catch(err => { console.log(err); });
}

// 使用了动态反射，所以每次获取的数据，都是扁平的一层数据，不需要再递归处理data本身
function putInDataList(keyMap, dataList, data) {
    let success = false;
    let rootKey = data.key.split(".")[0];
    if (!(rootKey in keyMap)) {
        dataList.push(data)
        success = true;
        return success;
    }

    for (let i = 0; i < dataList.length; i ++) {
        if (dataList[i].key == data.key) {
            dataList[i] = data;
            success = true;
        }
        else if (dataList[i].children != null) {
            success = putInDataList(keyMap, dataList[i].children, data);
        }
        if(success){
            return success;
        }
    }
}

// 将key数据放到keyMap中
function putInKeyMap(keyMap, data) {
    let keys = [];
    let names = [];
    keys.push(data.key);
    names.push(data.name);
    if (data.children != null) {
        for (let i = 0; i < data.children.length; i ++) {
            keys.push(data.children[i].key);
            names.push(data.children[i].name);
        }
    }
    for (let i = 0; i < keys.length; i ++) {
        if (!(keys[i] in keyMap)) {
            keyMap[keys[i]] = names[i];
        }
    }
}

export default {
    namespace: "reflectTree",

    state: {
        searchText: "",
        autoExpandParent: true,
        treeNodeData: [],
        expandedKeys: [],
        keyMap: {}, //所有节点的{key: title}哈希表，用于快速搜索
    },

    subscriptions: {
        setup({ dispatch, history }) {
            dispatch({
                type: "fetchData",
                payload: { key: "all" },
            });
        },
    },

    effects: {
        *search({ payload }, { call, put, select }) {
            let value = payload.data.target.value;
            yield put({
                type: "updateSearchText",
                payload: { data: value }
            });
            if (value == "") {
                return;
            }
            let expandedKeys = [];
            let keyMap = yield select(state => state.reflectTree.keyMap);
            for (let key in keyMap) {
                if (keyMap[key].toLowerCase().indexOf(value) != -1) {
                    expandedKeys.push(key);
                }
            }
            yield put({
                type: "updateExpandedKeys",
                payload: { data: expandedKeys }
            });
            yield put({
                type: "updateAutoExpandParent",
                payload: { data: true }
            })
        },
        *expand({ payload }, { call, put, select }) {
            let expandedKeys = payload.expandedKeys;
            let expanded = payload.expanded;
            let node = payload.node;
            yield put({
                type: "updateExpandedKeys",
                payload: { data: expandedKeys }
            });
            yield put({
                type: "updateAutoExpandParent",
                payload: { data: false }
            })
            if (expanded) {
                yield put({
                    type: "fetchData",
                    payload: { key: node.props.eventKey }
                });
            }
        },
        *fetchData({ payload }, { call, put, select }) {
            let data = yield call(getData, payload.key);
            if (data != null) {
                let dataList = yield select( state => state.reflectTree.treeNodeData);
                let keyMap = yield select( state => state.reflectTree.keyMap );
                for (let i = 0; i < data.length; i ++) {
                    putInDataList(keyMap, dataList, data[i]); // 将node数据放到列表中
                    putInKeyMap(keyMap, data[i]); // 将key数据放到keyMap中
                }
                yield put({
                    type: "updateKeyMap",
                    payload: {data: keyMap}
                });
                yield put({
                    type: "updateDataList",
                    payload: { data: dataList }
                });
            }
        },
    },

    reducers: {
        updateSearchText(state, { payload }) {
            return { ...state, searchText: payload.data }
        },
        updateAutoExpandParent(state, { payload }) {
            return { ...state, autoExpandParent: payload.data }
        },
        updateExpandedKeys(state, { payload }) {
            return { ...state, expandedKeys: [...payload.data] };
        },
        updateKeyMap(state, { payload }) {
            return { ...state, keyMap: {...payload.data} };
        },
        updateDataList(state, { payload }) {
            return { ...state, treeNodeData: [...payload.data] };
        },
    },
};