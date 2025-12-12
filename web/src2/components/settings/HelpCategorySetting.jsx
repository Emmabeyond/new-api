/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React, { useState, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Table,
  Button,
  Modal,
  Form,
  Tag,
  Space,
  Popconfirm,
  Empty,
  Spin,
  Typography,
  Card,
  SideSheet,
  Row,
  Col,
  Avatar,
} from '@douyinfe/semi-ui';
import {
  IconPlus,
  IconEdit,
  IconDelete,
  IconSave,
  IconClose,
  IconFolder,
} from '@douyinfe/semi-icons';
import {
  IllustrationNoResult,
  IllustrationNoResultDark,
} from '@douyinfe/semi-illustrations';
import { API, showError, showSuccess } from '../../helpers';
import { useIsMobile } from '../../hooks/common/useIsMobile';

const { Text, Title } = Typography;

const HelpCategorySetting = () => {
  const { t } = useTranslation();
  const isMobile = useIsMobile();
  const formApiRef = useRef(null);

  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showEdit, setShowEdit] = useState(false);
  const [editingCat, setEditingCat] = useState({});
  const [submitting, setSubmitting] = useState(false);

  const isEdit = editingCat.id !== undefined;

  useEffect(() => {
    loadCategories();
  }, []);

  const loadCategories = async () => {
    setLoading(true);
    try {
      const res = await API.get('/api/help/admin/categories');
      if (res.data.success) {
        setCategories(res.data.data || []);
      } else {
        showError(res.data.message || t('加载分类失败'));
      }
    } catch (error) {
      showError(t('加载分类失败'));
    } finally {
      setLoading(false);
    }
  };

  const getInitValues = () => ({
    name: '',
    sort_order: 0,
    status: 1,
  });

  const handleCreate = () => {
    setEditingCat({});
    setShowEdit(true);
  };

  const handleEdit = (record) => {
    setEditingCat(record);
    setShowEdit(true);
  };

  const handleClose = () => {
    setShowEdit(false);
    setEditingCat({});
  };

  const handleDelete = async (id, documentCount) => {
    if (documentCount > 0) {
      Modal.warning({
        title: t('无法删除'),
        content: t('该分类下还有 {{count}} 篇文档，请先删除或移动这些文档', { count: documentCount }),
      });
      return;
    }

    try {
      const res = await API.delete(`/api/help/admin/categories/${id}`);
      if (res.data.success) {
        showSuccess(t('删除成功'));
        loadCategories();
      } else {
        showError(res.data.message || t('删除失败'));
      }
    } catch (error) {
      showError(t('删除失败'));
    }
  };

  const handleSubmit = async (values) => {
    if (!values.name?.trim()) {
      showError(t('请输入分类名称'));
      return;
    }

    setSubmitting(true);
    try {
      const payload = {
        ...values,
        sort_order: values.sort_order || 0,
        status: values.status ? 1 : 2,
      };

      let res;
      if (isEdit) {
        res = await API.put(`/api/help/admin/categories/${editingCat.id}`, payload);
      } else {
        res = await API.post('/api/help/admin/categories', payload);
      }

      if (res.data.success) {
        showSuccess(isEdit ? t('更新成功') : t('创建成功'));
        handleClose();
        loadCategories();
      } else {
        showError(res.data.message || t('操作失败'));
      }
    } catch (error) {
      showError(t('操作失败'));
    } finally {
      setSubmitting(false);
    }
  };

  const columns = [
    {
      title: t('分类名称'),
      dataIndex: 'name',
      key: 'name',
      render: (text) => <Text strong>{text}</Text>,
    },
    {
      title: t('排序'),
      dataIndex: 'sort_order',
      key: 'sort_order',
      width: 80,
    },
    {
      title: t('文档数量'),
      dataIndex: 'document_count',
      key: 'document_count',
      width: 100,
      render: (count) => (
        <Tag color='blue'>{count || 0}</Tag>
      ),
    },
    {
      title: t('状态'),
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (status) => (
        <Tag color={status === 1 ? 'green' : 'grey'}>
          {status === 1 ? t('启用') : t('禁用')}
        </Tag>
      ),
    },
    {
      title: t('操作'),
      key: 'operate',
      width: 150,
      render: (_, record) => (
        <Space>
          <Button
            theme='light'
            type='primary'
            icon={<IconEdit />}
            size='small'
            onClick={() => handleEdit(record)}
          />
          <Popconfirm
            title={
              record.document_count > 0
                ? t('该分类下有文档，无法删除')
                : t('确定删除该分类吗？')
            }
            onConfirm={() => handleDelete(record.id, record.document_count)}
            disabled={record.document_count > 0}
          >
            <Button
              theme='light'
              type='danger'
              icon={<IconDelete />}
              size='small'
              disabled={record.document_count > 0}
            />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  useEffect(() => {
    if (showEdit && formApiRef.current) {
      if (isEdit) {
        formApiRef.current.setValues({
          name: editingCat.name || '',
          sort_order: editingCat.sort_order || 0,
          status: editingCat.status === 1,
        });
      } else {
        formApiRef.current.setValues({
          ...getInitValues(),
          status: true,
        });
      }
    }
  }, [showEdit, editingCat, isEdit]);

  return (
    <div className='p-4'>
      <Card className='!rounded-2xl'>
        <div className='flex justify-between items-center mb-4'>
          <div className='flex items-center gap-2'>
            <IconFolder size='large' />
            <Title heading={5} className='!m-0'>
              {t('帮助分类管理')}
            </Title>
          </div>
          <Button
            theme='solid'
            type='primary'
            icon={<IconPlus />}
            onClick={handleCreate}
          >
            {t('新建分类')}
          </Button>
        </div>

        <Table
          columns={columns}
          dataSource={categories}
          loading={loading}
          rowKey='id'
          pagination={{
            pageSize: 10,
            showSizeChanger: true,
            pageSizeOptions: [10, 20, 50],
          }}
          empty={
            <Empty
              image={<IllustrationNoResult style={{ width: 150, height: 150 }} />}
              darkModeImage={<IllustrationNoResultDark style={{ width: 150, height: 150 }} />}
              description={t('暂无分类')}
            />
          }
        />
      </Card>

      <SideSheet
        placement='right'
        title={
          <Space>
            {isEdit ? (
              <Tag color='blue' shape='circle'>{t('编辑')}</Tag>
            ) : (
              <Tag color='green' shape='circle'>{t('新建')}</Tag>
            )}
            <Title heading={4} className='!m-0'>
              {isEdit ? t('编辑分类') : t('新建分类')}
            </Title>
          </Space>
        }
        visible={showEdit}
        width={isMobile ? '100%' : 480}
        onCancel={handleClose}
        footer={
          <div className='flex justify-end'>
            <Space>
              <Button
                theme='solid'
                type='primary'
                icon={<IconSave />}
                loading={submitting}
                onClick={() => formApiRef.current?.submitForm()}
              >
                {t('保存')}
              </Button>
              <Button
                theme='light'
                type='tertiary'
                icon={<IconClose />}
                onClick={handleClose}
              >
                {t('取消')}
              </Button>
            </Space>
          </div>
        }
      >
        <Spin spinning={submitting}>
          <Form
            getFormApi={(api) => (formApiRef.current = api)}
            onSubmit={handleSubmit}
            initValues={getInitValues()}
          >
            <Card className='!rounded-2xl'>
              <div className='flex items-center mb-4'>
                <Avatar size='small' color='blue' className='mr-2'>
                  <IconFolder size={16} />
                </Avatar>
                <div>
                  <Text className='text-lg font-medium'>{t('分类信息')}</Text>
                  <div className='text-xs text-gray-500'>{t('设置分类的基本信息')}</div>
                </div>
              </div>

              <Row gutter={16}>
                <Col span={24}>
                  <Form.Input
                    field='name'
                    label={t('分类名称')}
                    placeholder={t('请输入分类名称')}
                    rules={[{ required: true, message: t('请输入分类名称') }]}
                  />
                </Col>
                <Col span={12}>
                  <Form.InputNumber
                    field='sort_order'
                    label={t('排序权重')}
                    placeholder='0'
                    min={0}
                  />
                </Col>
                <Col span={12}>
                  <Form.Switch
                    field='status'
                    label={t('启用状态')}
                  />
                </Col>
              </Row>
            </Card>
          </Form>
        </Spin>
      </SideSheet>
    </div>
  );
};

export default HelpCategorySetting;
