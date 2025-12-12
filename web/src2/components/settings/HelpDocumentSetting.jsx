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
  Input,
  Select,
  InputNumber,
  Switch,
  SideSheet,
  Row,
  Col,
  Avatar,
  TextArea,
} from '@douyinfe/semi-ui';
import {
  IconPlus,
  IconEdit,
  IconDelete,
  IconSave,
  IconClose,
  IconFile,
  IconList,
} from '@douyinfe/semi-icons';
import {
  IllustrationNoResult,
  IllustrationNoResultDark,
} from '@douyinfe/semi-illustrations';
import { API, showError, showSuccess } from '../../helpers';
import { useIsMobile } from '../../hooks/common/useIsMobile';
import MarkdownEditor from './MarkdownEditor';

const { Text, Title } = Typography;

const HelpDocumentSetting = () => {
  const { t } = useTranslation();
  const isMobile = useIsMobile();
  const formApiRef = useRef(null);

  const [documents, setDocuments] = useState([]);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showEdit, setShowEdit] = useState(false);
  const [editingDoc, setEditingDoc] = useState({});
  const [submitting, setSubmitting] = useState(false);
  const [content, setContent] = useState('');

  const isEdit = editingDoc.id !== undefined;

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [docsRes, catsRes] = await Promise.all([
        API.get('/api/help/admin/documents'),
        API.get('/api/help/admin/categories'),
      ]);

      if (docsRes.data.success) {
        setDocuments(docsRes.data.data || []);
      } else {
        showError(docsRes.data.message || t('加载文档失败'));
      }

      if (catsRes.data.success) {
        setCategories(catsRes.data.data || []);
      }
    } catch (error) {
      showError(t('加载数据失败'));
    } finally {
      setLoading(false);
    }
  };

  const getInitValues = () => ({
    title: '',
    content: '',
    category_id: 0,
    sort_order: 0,
    status: 1,
  });

  const handleCreate = () => {
    setEditingDoc({});
    setShowEdit(true);
  };

  const handleEdit = async (record) => {
    setEditingDoc(record);
    setShowEdit(true);
  };

  const handleClose = () => {
    setShowEdit(false);
    setEditingDoc({});
  };

  const handleDelete = async (id) => {
    try {
      const res = await API.delete(`/api/help/admin/documents/${id}`);
      if (res.data.success) {
        showSuccess(t('删除成功'));
        loadData();
      } else {
        showError(res.data.message || t('删除失败'));
      }
    } catch (error) {
      showError(t('删除失败'));
    }
  };

  const handleSubmit = async (values) => {
    if (!values.title?.trim()) {
      showError(t('请输入文档标题'));
      return;
    }
    if (!content?.trim()) {
      showError(t('请输入文档内容'));
      return;
    }

    setSubmitting(true);
    try {
      const payload = {
        ...values,
        content,
        category_id: values.category_id || 0,
        sort_order: values.sort_order || 0,
        status: values.status ? 1 : 2,
      };

      let res;
      if (isEdit) {
        res = await API.put(`/api/help/admin/documents/${editingDoc.id}`, payload);
      } else {
        res = await API.post('/api/help/admin/documents', payload);
      }

      if (res.data.success) {
        showSuccess(isEdit ? t('更新成功') : t('创建成功'));
        handleClose();
        loadData();
      } else {
        showError(res.data.message || t('操作失败'));
      }
    } catch (error) {
      showError(t('操作失败'));
    } finally {
      setSubmitting(false);
    }
  };

  const getCategoryName = (categoryId) => {
    const cat = categories.find((c) => c.id === categoryId);
    return cat ? cat.name : t('未分类');
  };

  const columns = [
    {
      title: t('标题'),
      dataIndex: 'title',
      key: 'title',
      render: (text) => <Text strong>{text}</Text>,
    },
    {
      title: t('分类'),
      dataIndex: 'category_id',
      key: 'category_id',
      render: (categoryId) => (
        <Tag color='blue'>{getCategoryName(categoryId)}</Tag>
      ),
    },
    {
      title: t('排序'),
      dataIndex: 'sort_order',
      key: 'sort_order',
      width: 80,
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
            title={t('确定删除该文档吗？')}
            onConfirm={() => handleDelete(record.id)}
          >
            <Button
              theme='light'
              type='danger'
              icon={<IconDelete />}
              size='small'
            />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const handleAfterVisibleChange = (visible) => {
    if (visible && formApiRef.current) {
      if (isEdit) {
        formApiRef.current.setValues({
          title: editingDoc.title || '',
          category_id: editingDoc.category_id || 0,
          sort_order: editingDoc.sort_order || 0,
          status: editingDoc.status === 1,
        });
        setContent(editingDoc.content || '');
      } else {
        formApiRef.current.setValues({
          ...getInitValues(),
          status: true,
        });
        setContent('');
      }
    }
  };

  return (
    <div className='p-4'>
      <Card className='!rounded-2xl'>
        <div className='flex justify-between items-center mb-4'>
          <div className='flex items-center gap-2'>
            <IconFile size='large' />
            <Title heading={5} className='!m-0'>
              {t('帮助文档管理')}
            </Title>
          </div>
          <Button
            theme='solid'
            type='primary'
            icon={<IconPlus />}
            onClick={handleCreate}
          >
            {t('新建文档')}
          </Button>
        </div>

        <Table
          columns={columns}
          dataSource={documents}
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
              description={t('暂无文档')}
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
              {isEdit ? t('编辑文档') : t('新建文档')}
            </Title>
          </Space>
        }
        visible={showEdit}
        width={isMobile ? '100%' : 800}
        onCancel={handleClose}
        afterVisibleChange={handleAfterVisibleChange}
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
            <Card className='!rounded-2xl mb-4'>
              <div className='flex items-center mb-4'>
                <Avatar size='small' color='blue' className='mr-2'>
                  <IconFile size={16} />
                </Avatar>
                <div>
                  <Text className='text-lg font-medium'>{t('基本信息')}</Text>
                  <div className='text-xs text-gray-500'>{t('设置文档的基本信息')}</div>
                </div>
              </div>

              <Row gutter={16}>
                <Col span={24}>
                  <Form.Input
                    field='title'
                    label={t('文档标题')}
                    placeholder={t('请输入文档标题')}
                    rules={[{ required: true, message: t('请输入文档标题') }]}
                  />
                </Col>
                <Col span={12}>
                  <Form.Select
                    field='category_id'
                    label={t('所属分类')}
                    placeholder={t('请选择分类')}
                    optionList={[
                      { value: 0, label: t('未分类') },
                      ...categories.map((cat) => ({
                        value: cat.id,
                        label: cat.name,
                      })),
                    ]}
                  />
                </Col>
                <Col span={6}>
                  <Form.InputNumber
                    field='sort_order'
                    label={t('排序权重')}
                    placeholder='0'
                    min={0}
                  />
                </Col>
                <Col span={6}>
                  <Form.Switch
                    field='status'
                    label={t('启用状态')}
                  />
                </Col>
              </Row>
            </Card>

            <Card className='!rounded-2xl'>
              <div className='flex items-center mb-4'>
                <Avatar size='small' color='green' className='mr-2'>
                  <IconList size={16} />
                </Avatar>
                <div>
                  <Text className='text-lg font-medium'>{t('文档内容')}</Text>
                  <div className='text-xs text-gray-500'>{t('使用 Markdown 格式编写文档内容')}</div>
                </div>
              </div>

              <MarkdownEditor
                value={content}
                onChange={setContent}
                height={400}
              />
            </Card>
          </Form>
        </Spin>
      </SideSheet>
    </div>
  );
};

export default HelpDocumentSetting;
