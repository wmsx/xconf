import React, { FC } from 'react';
import moment from 'moment';
import { Button, Col, Form, Input, Modal, Row } from 'antd';
import { ModalProps } from 'antd/lib/modal';

import { withModal } from '@src/hoc';
import { useCapture } from '@src/hooks';
import { ReleaseForm, releaseConfig } from '@src/services';

export interface ReleaseModelProps extends Omit<ModalProps, 'onOk' | 'onCancel'> {
  appName: string;
  clusterName: string;
  namespaceName: string;
  onOk?: () => void;
  onCancel?: () => void;
}

const ReleaseModel: FC<ReleaseModelProps> = ({ appName, clusterName, namespaceName, onOk, ...props }) => {
  const [form] = Form.useForm();
  const [releaseState, release] = useCapture<any, ReleaseForm>({
    fn: releaseConfig,
    initialState: {},
    option: {
      onSuccess: () => {
        onOk && onOk();
        props.onCancel && props.onCancel();
      },
    },
  });

  return (
    <Modal {...props} footer={null} title={`${appName}/${clusterName}/${namespaceName}`}>
      <Form
        name="release_form"
        form={form}
        labelCol={{ span: 4 }}
        wrapperCol={{ span: 20 }}
        initialValues={{ tag: 'v' + moment().format('YYYY-MM-DD') }}
        onFinish={({ tag, comment }) => release({ appName, clusterName, namespaceName, tag, comment })}
      >
        <Form.Item name="tag" label="tag" rules={[{ required: true, message: '必须填写tag' }]}>
          <Input placeholder="请输入Tag" />
        </Form.Item>
        <Form.Item name="comment" label="备注">
          <Input.TextArea placeholder="发布备注" />
        </Form.Item>
        <Row gutter={16}>
          <Col offset={4} span={10}>
            <Button type="primary" htmlType="submit" loading={releaseState.loading} style={{ width: '100%' }}>
              发布
            </Button>
          </Col>
          <Col span={10}>
            <Button htmlType="button" onClick={props.onCancel} style={{ width: '100%' }}>
              取消
            </Button>
          </Col>
        </Row>
      </Form>
    </Modal>
  );
};

export default withModal<ReleaseModelProps>(ReleaseModel);
