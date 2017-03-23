import React from 'react';
import { connect } from 'dva';
import styles from './IndexPage.css';

import ReflectTree from "../components/ReflectTree"; 

function IndexPage(props) {
  return (
    <div>
      <ReflectTree />
    </div>
  );
}

IndexPage.propTypes = {
};

export default connect()(IndexPage);
