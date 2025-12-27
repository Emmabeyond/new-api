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

import React, { useEffect, useState, useContext, useRef, useCallback } from 'react';
import {
  API,
  showError,
  showInfo,
  showSuccess,
  renderQuota,
  renderQuotaWithAmount,
  copy,
  getQuotaPerUnit,
} from '../../helpers';
import { Modal, Toast } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';
import { UserContext } from '../../context/User';
import { StatusContext } from '../../context/Status';

import RechargeCard from './RechargeCard';
import InvitationCard from './InvitationCard';
import TransferModal from './modals/TransferModal';
import PaymentConfirmModal from './modals/PaymentConfirmModal';
import TopupHistoryModal from './modals/TopupHistoryModal';
import CheckinResultModal from './modals/CheckinResultModal';
import QRCodePaymentModal from './modals/QRCodePaymentModal';
import { SliderCaptcha } from '../captcha';

const TopUp = () => {
  const { t } = useTranslation();
  const [userState, userDispatch] = useContext(UserContext);
  const [statusState] = useContext(StatusContext);

  const [redemptionCode, setRedemptionCode] = useState('');
  const [amount, setAmount] = useState(0.0);
  const [minTopUp, setMinTopUp] = useState(statusState?.status?.min_topup || 1);
  const [topUpCount, setTopUpCount] = useState(
    statusState?.status?.min_topup || 1,
  );
  const [topUpLink, setTopUpLink] = useState(
    statusState?.status?.top_up_link || '',
  );
  const [enableOnlineTopUp, setEnableOnlineTopUp] = useState(
    statusState?.status?.enable_online_topup || false,
  );
  const [priceRatio, setPriceRatio] = useState(statusState?.status?.price || 1);

  const [enableStripeTopUp, setEnableStripeTopUp] = useState(
    statusState?.status?.enable_stripe_topup || false,
  );
  const [statusLoading, setStatusLoading] = useState(true);

  // Creem 相关状态
  const [creemProducts, setCreemProducts] = useState([]);
  const [enableCreemTopUp, setEnableCreemTopUp] = useState(false);
  const [creemOpen, setCreemOpen] = useState(false);
  const [selectedCreemProduct, setSelectedCreemProduct] = useState(null);

  // LINUX DO Credit 相关状态
  const [enableLinuxDoTopUp, setEnableLinuxDoTopUp] = useState(false);
  const [linuxDoMinTopUp, setLinuxDoMinTopUp] = useState(1);
  const [linuxDoLoading, setLinuxDoLoading] = useState(false);

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [open, setOpen] = useState(false);
  const [payWay, setPayWay] = useState('');
  const [amountLoading, setAmountLoading] = useState(false);
  const [paymentLoading, setPaymentLoading] = useState(false);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [payMethods, setPayMethods] = useState([]);

  const affFetchedRef = useRef(false);

  // 邀请相关状态
  const [affLink, setAffLink] = useState('');
  const [openTransfer, setOpenTransfer] = useState(false);
  const [transferAmount, setTransferAmount] = useState(0);

  // 账单Modal状态
  const [openHistory, setOpenHistory] = useState(false);

  // 签到相关状态
  const [checkinStats, setCheckinStats] = useState(null);
  const [checkinLoading, setCheckinLoading] = useState(false);
  const [sliderCaptchaEnabled, setSliderCaptchaEnabled] = useState(false);
  const [captchaModalVisible, setCaptchaModalVisible] = useState(false);
  const [checkinResultVisible, setCheckinResultVisible] = useState(false);
  const [checkinResult, setCheckinResult] = useState(null);

  // 预设充值额度选项
  const [presetAmounts, setPresetAmounts] = useState([]);
  const [selectedPreset, setSelectedPreset] = useState(null);

  // 充值配置信息
  const [topupInfo, setTopupInfo] = useState({
    amount_options: [],
    discount: {},
  });

  // 二维码支付弹窗状态
  const [qrCodeModalVisible, setQrCodeModalVisible] = useState(false);
  const [qrCodeData, setQrCodeData] = useState({
    qrCodeUrl: '',
    tradeNo: '',
    amount: 0,
    money: '',
    paymentMethod: '',
  });
  const [qrCodeLoading, setQrCodeLoading] = useState(false);

  const topUp = async () => {
    if (redemptionCode === '') {
      showInfo(t('请输入兑换码！'));
      return;
    }
    setIsSubmitting(true);
    try {
      const res = await API.post('/api/user/topup', {
        key: redemptionCode,
      });
      const { success, message, data } = res.data;
      if (success) {
        showSuccess(t('兑换成功！'));
        Modal.success({
          title: t('兑换成功！'),
          content: t('成功兑换额度：') + renderQuota(data),
          centered: true,
        });
        if (userState.user) {
          const updatedUser = {
            ...userState.user,
            quota: userState.user.quota + data,
          };
          userDispatch({ type: 'login', payload: updatedUser });
        }
        setRedemptionCode('');
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('请求失败'));
    } finally {
      setIsSubmitting(false);
    }
  };

  const openTopUpLink = () => {
    if (!topUpLink) {
      showError(t('超级管理员未设置充值链接！'));
      return;
    }
    window.open(topUpLink, '_blank');
  };

  const preTopUp = async (payment) => {
    if (payment === 'stripe') {
      if (!enableStripeTopUp) {
        showError(t('管理员未开启Stripe充值！'));
        return;
      }
    } else if (payment === 'linuxdo') {
      if (!enableLinuxDoTopUp) {
        showError(t('管理员未开启 LINUX DO Credit 充值！'));
        return;
      }
    } else {
      if (!enableOnlineTopUp) {
        showError(t('管理员未开启在线充值！'));
        return;
      }
    }

    if (topUpCount < minTopUp) {
      showError(t('充值数量不能小于') + minTopUp);
      return;
    }

    setPayWay(payment);
    setPaymentLoading(true);

    try {
      // LINUX DO Credit: 直接跳转支付页面，不需要二维码
      if (payment === 'linuxdo') {
        await requestLinuxDoPayment();
        return;
      }

      // 对于支付宝、微信支付，尝试使用二维码模式
      if (payment === 'alipay' || payment === 'wxpay') {
        const qrCodeResult = await requestQRCode(payment);
        if (qrCodeResult) {
          // 成功获取二维码，显示二维码弹窗
          return;
        }
        // 获取二维码失败，降级到跳转模式
        showInfo(t('正在跳转到支付页面...'));
      }

      // Stripe 或降级模式：使用原有的确认弹窗流程
      if (payment === 'stripe') {
        await getStripeAmount();
      } else {
        await getAmount();
        setOpen(true);
      }
    } catch (error) {
      showError(t('获取金额失败'));
    } finally {
      setPaymentLoading(false);
    }
  };

  // 请求 LINUX DO Credit 支付（直接跳转）
  const requestLinuxDoPayment = async () => {
    try {
      const res = await API.post('/api/user/pay/qrcode', {
        amount: parseInt(topUpCount),
        payment_method: 'linuxdo',
      });

      if (res.data.message === 'success') {
        const data = res.data.data;
        if (data.payurl) {
          // 直接跳转到支付页面
          window.open(data.payurl, '_blank');
          return true;
        }
      } else {
        showError(res.data.data || t('获取支付链接失败'));
      }
      return false;
    } catch (err) {
      console.error('LINUX DO Credit 支付请求失败:', err);
      return false;
    } finally {
      setPaymentLoading(false);
    }
  };

  // 请求二维码支付
  const requestQRCode = async (payment) => {
    try {
      setQrCodeLoading(true);
      const res = await API.post('/api/user/pay/qrcode', {
        amount: parseInt(topUpCount),
        payment_method: payment,
      });

      if (res.data.message === 'success') {
        const data = res.data.data;
        
        // 优先使用二维码，其次使用支付链接
        if (data.qrcode) {
          setQrCodeData({
            qrCodeUrl: data.qrcode,
            tradeNo: data.trade_no,
            amount: data.amount,
            money: data.money,
            paymentMethod: payment,
          });
          setQrCodeModalVisible(true);
          setPaymentLoading(false);
          return true;
        } else if (data.payurl) {
          // 没有二维码但有支付链接，降级到跳转模式
          window.open(data.payurl, '_blank');
          setPaymentLoading(false);
          return true;
        }
      } else {
        // API 返回错误
        showError(res.data.data || t('获取支付二维码失败'));
      }
      return false;
    } catch (err) {
      console.error('请求二维码失败:', err);
      return false;
    } finally {
      setQrCodeLoading(false);
    }
  };

  // 二维码支付成功回调
  const handleQRCodePaymentSuccess = (rechargedAmount) => {
    setQrCodeModalVisible(false);
    
    // 显示成功提示
    showSuccess(t('充值成功！'));
    
    // 计算实际充值的额度
    const quotaPerUnit = getQuotaPerUnit() || 500000;
    const quotaToAdd = rechargedAmount * quotaPerUnit;
    
    Modal.success({
      title: t('充值成功！'),
      content: t('成功充值额度：') + renderQuota(quotaToAdd),
      centered: true,
    });

    // 更新用户额度
    if (userState.user) {
      const updatedUser = {
        ...userState.user,
        quota: userState.user.quota + quotaToAdd,
      };
      userDispatch({ type: 'login', payload: updatedUser });
    }
  };

  // 关闭二维码弹窗
  const handleQRCodeModalClose = () => {
    setQrCodeModalVisible(false);
    setQrCodeData({
      qrCodeUrl: '',
      tradeNo: '',
      amount: 0,
      money: '',
      paymentMethod: '',
    });
  };

  // 刷新二维码
  const handleQRCodeRefresh = async () => {
    const payment = qrCodeData.paymentMethod;
    if (payment) {
      await requestQRCode(payment);
    }
  };

  const onlineTopUp = async () => {
    if (payWay === 'stripe') {
      // Stripe 支付处理
      if (amount === 0) {
        await getStripeAmount();
      }
    } else {
      // 普通支付处理
      if (amount === 0) {
        await getAmount();
      }
    }

    if (topUpCount < minTopUp) {
      showError('充值数量不能小于' + minTopUp);
      return;
    }
    setConfirmLoading(true);
    try {
      let res;
      if (payWay === 'stripe') {
        // Stripe 支付请求
        res = await API.post('/api/user/stripe/pay', {
          amount: parseInt(topUpCount),
          payment_method: 'stripe',
        });
      } else {
        // 普通支付请求
        res = await API.post('/api/user/pay', {
          amount: parseInt(topUpCount),
          payment_method: payWay,
        });
      }

      if (res !== undefined) {
        const { message, data } = res.data;
        if (message === 'success') {
          if (payWay === 'stripe') {
            // Stripe 支付回调处理
            window.open(data.pay_link, '_blank');
          } else {
            // 普通支付表单提交
            let params = data;
            let url = res.data.url;
            let form = document.createElement('form');
            form.action = url;
            form.method = 'POST';
            let isSafari =
              navigator.userAgent.indexOf('Safari') > -1 &&
              navigator.userAgent.indexOf('Chrome') < 1;
            if (!isSafari) {
              form.target = '_blank';
            }
            for (let key in params) {
              let input = document.createElement('input');
              input.type = 'hidden';
              input.name = key;
              input.value = params[key];
              form.appendChild(input);
            }
            document.body.appendChild(form);
            form.submit();
            document.body.removeChild(form);
          }
        } else {
          showError(data);
        }
      } else {
        showError(res);
      }
    } catch (err) {
      showError(t('支付请求失败'));
    } finally {
      setOpen(false);
      setConfirmLoading(false);
    }
  };

  const creemPreTopUp = async (product) => {
    if (!enableCreemTopUp) {
      showError(t('管理员未开启 Creem 充值！'));
      return;
    }
    setSelectedCreemProduct(product);
    setCreemOpen(true);
  };

  const onlineCreemTopUp = async () => {
    if (!selectedCreemProduct) {
      showError(t('请选择产品'));
      return;
    }
    // Validate product has required fields
    if (!selectedCreemProduct.productId) {
      showError(t('产品配置错误，请联系管理员'));
      return;
    }
    setConfirmLoading(true);
    try {
      const res = await API.post('/api/user/creem/pay', {
        product_id: selectedCreemProduct.productId,
        payment_method: 'creem',
      });
      if (res !== undefined) {
        const { message, data } = res.data;
        if (message === 'success') {
          processCreemCallback(data);
        } else {
          showError(data);
        }
      } else {
        showError(res);
      }
    } catch (err) {
      showError(t('支付请求失败'));
    } finally {
      setCreemOpen(false);
      setConfirmLoading(false);
    }
  };

  const processCreemCallback = (data) => {
    // 与 Stripe 保持一致的实现方式
    window.open(data.checkout_url, '_blank');
  };


  const getUserQuota = async () => {
    let res = await API.get(`/api/user/self`);
    const { success, message, data } = res.data;
    if (success) {
      userDispatch({ type: 'login', payload: data });
    } else {
      showError(message);
    }
  };

  // 获取充值配置信息
  const getTopupInfo = async () => {
    try {
      const res = await API.get('/api/user/topup/info');
      const { message, data, success } = res.data;
      if (success) {
        setTopupInfo({
          amount_options: data.amount_options || [],
          discount: data.discount || {},
        });

        // 处理支付方式
        let payMethods = data.pay_methods || [];
        try {
          if (typeof payMethods === 'string') {
            payMethods = JSON.parse(payMethods);
          }
          if (payMethods && payMethods.length > 0) {
            // 检查name和type是否为空
            payMethods = payMethods.filter((method) => {
              return method.name && method.type;
            });
            // 如果没有color，则设置默认颜色
            payMethods = payMethods.map((method) => {
              // 规范化最小充值数
              const normalizedMinTopup = Number(method.min_topup);
              method.min_topup = Number.isFinite(normalizedMinTopup)
                ? normalizedMinTopup
                : 0;

              // Stripe 的最小充值从后端字段回填
              if (
                method.type === 'stripe' &&
                (!method.min_topup || method.min_topup <= 0)
              ) {
                const stripeMin = Number(data.stripe_min_topup);
                if (Number.isFinite(stripeMin)) {
                  method.min_topup = stripeMin;
                }
              }

              if (!method.color) {
                if (method.type === 'alipay') {
                  method.color = 'rgba(var(--semi-blue-5), 1)';
                } else if (method.type === 'wxpay') {
                  method.color = 'rgba(var(--semi-green-5), 1)';
                } else if (method.type === 'stripe') {
                  method.color = 'rgba(var(--semi-purple-5), 1)';
                } else {
                  method.color = 'rgba(var(--semi-primary-5), 1)';
                }
              }
              return method;
            });
          } else {
            payMethods = [];
          }

          // 如果启用了 Stripe 支付，添加到支付方法列表
          // 这个逻辑现在由后端处理，如果 Stripe 启用，后端会在 pay_methods 中包含它

          setPayMethods(payMethods);
          const enableStripeTopUp = data.enable_stripe_topup || false;
          const enableOnlineTopUp = data.enable_online_topup || false;
          const enableCreemTopUp = data.enable_creem_topup || false;
          const minTopUpValue = enableOnlineTopUp
            ? data.min_topup
            : enableStripeTopUp
              ? data.stripe_min_topup
              : 1;
          const enableLinuxDoTopUp = data.enable_linuxdo_topup || false;
          setEnableOnlineTopUp(enableOnlineTopUp);
          setEnableStripeTopUp(enableStripeTopUp);
          setEnableCreemTopUp(enableCreemTopUp);
          setEnableLinuxDoTopUp(enableLinuxDoTopUp);
          setLinuxDoMinTopUp(data.linuxdo_min_topup || 1);
          setMinTopUp(minTopUpValue);
          setTopUpCount(minTopUpValue);

          // 设置 Creem 产品
          try {
            const products = JSON.parse(data.creem_products || '[]');
            setCreemProducts(products);
          } catch (e) {
            setCreemProducts([]);
          }

          // 如果没有自定义充值数量选项，根据最小充值金额生成预设充值额度选项
          if (topupInfo.amount_options.length === 0) {
            setPresetAmounts(generatePresetAmounts(minTopUpValue));
          }

          // 初始化显示实付金额
          getAmount(minTopUpValue);
        } catch (e) {
          setPayMethods([]);
        }

        // 如果有自定义充值数量选项，使用它们替换默认的预设选项
        if (data.amount_options && data.amount_options.length > 0) {
          const customPresets = data.amount_options.map((amount) => ({
            value: amount,
            discount: data.discount[amount] || 1.0,
          }));
          setPresetAmounts(customPresets);
        }
      } else {
        console.error('获取充值配置失败:', data);
      }
    } catch (error) {
      console.error('获取充值配置异常:', error);
    }
  };

  // 获取邀请链接
  const getAffLink = async () => {
    const res = await API.get('/api/user/aff');
    const { success, message, data } = res.data;
    if (success) {
      let link = `${window.location.origin}/register?aff=${data}`;
      setAffLink(link);
    } else {
      showError(message);
    }
  };

  // 划转邀请额度
  const transfer = async () => {
    if (transferAmount < getQuotaPerUnit()) {
      showError(t('划转金额最低为') + ' ' + renderQuota(getQuotaPerUnit()));
      return;
    }
    const res = await API.post(`/api/user/aff_transfer`, {
      quota: transferAmount,
    });
    const { success, message } = res.data;
    if (success) {
      showSuccess(message);
      setOpenTransfer(false);
      getUserQuota().then();
    } else {
      showError(message);
    }
  };

  // 复制邀请链接
  const handleAffLinkClick = async () => {
    await copy(affLink);
    showSuccess(t('邀请链接已复制到剪切板'));
  };

  useEffect(() => {
    if (!userState?.user?.id) {
      getUserQuota().then();
    }
    setTransferAmount(getQuotaPerUnit());
  }, []);

  useEffect(() => {
    if (affFetchedRef.current) return;
    affFetchedRef.current = true;
    getAffLink().then();
  }, []);

  // 在 statusState 可用时获取充值信息
  useEffect(() => {
    getTopupInfo().then();
  }, []);

  useEffect(() => {
    if (statusState?.status) {
      // const minTopUpValue = statusState.status.min_topup || 1;
      // setMinTopUp(minTopUpValue);
      // setTopUpCount(minTopUpValue);
      setTopUpLink(statusState.status.top_up_link || '');
      setPriceRatio(statusState.status.price || 1);

      setStatusLoading(false);
    }
  }, [statusState?.status]);

  const renderAmount = () => {
    return amount + ' ' + t('元');
  };

  const getAmount = async (value) => {
    if (value === undefined) {
      value = topUpCount;
    }
    setAmountLoading(true);
    try {
      const res = await API.post('/api/user/amount', {
        amount: parseFloat(value),
      });
      if (res !== undefined) {
        const { message, data } = res.data;
        if (message === 'success') {
          setAmount(parseFloat(data));
        } else {
          setAmount(0);
          Toast.error({ content: '错误：' + data, id: 'getAmount' });
        }
      } else {
        showError(res);
      }
    } catch (err) {
      // 忽略错误
    }
    setAmountLoading(false);
  };

  const getStripeAmount = async (value) => {
    if (value === undefined) {
      value = topUpCount;
    }
    setAmountLoading(true);
    try {
      const res = await API.post('/api/user/stripe/amount', {
        amount: parseFloat(value),
      });
      if (res !== undefined) {
        const { message, data } = res.data;
        if (message === 'success') {
          setAmount(parseFloat(data));
        } else {
          setAmount(0);
          Toast.error({ content: '错误：' + data, id: 'getAmount' });
        }
      } else {
        showError(res);
      }
    } catch (err) {
      // 忽略错误
    } finally {
      setAmountLoading(false);
    }
  };

  const handleCancel = () => {
    setOpen(false);
  };

  const handleTransferCancel = () => {
    setOpenTransfer(false);
  };

  const handleOpenHistory = () => {
    setOpenHistory(true);
  };

  const handleHistoryCancel = () => {
    setOpenHistory(false);
  };

  // 获取签到状态
  const fetchCheckinStats = useCallback(async () => {
    try {
      const res = await API.get('/api/user/checkin/stats');
      const { success, data, message } = res.data;
      if (success) {
        setCheckinStats(data);
      }
    } catch (err) {
      // 忽略错误，按钮保持默认可点击状态
    }
  }, []);

  // 获取验证码状态
  const fetchCaptchaStatus = useCallback(async () => {
    try {
      const res = await API.get('/api/captcha/status');
      if (res.data.success && res.data.data) {
        const { enabled, require_on_checkin } = res.data.data;
        setSliderCaptchaEnabled(enabled && require_on_checkin);
      }
    } catch (err) {
      // 忽略错误，默认不启用验证码
    }
  }, []);

  // 获取签到状态和验证码状态
  useEffect(() => {
    fetchCheckinStats();
    fetchCaptchaStatus();
  }, [fetchCheckinStats, fetchCaptchaStatus]);

  // 滑块验证码成功回调
  const handleSliderCaptchaSuccess = (token) => {
    setCaptchaModalVisible(false);
    doCheckin(token);
  };

  // 实际执行签到
  const doCheckin = async (captchaToken) => {
    setCheckinLoading(true);
    try {
      const checkinData = {};
      if (captchaToken) {
        checkinData.captcha_token = captchaToken;
      }
      const res = await API.post('/api/user/checkin', checkinData);
      const { success, data, message } = res.data;
      if (success) {
        setCheckinResult(data);
        setCheckinResultVisible(true);
        // 刷新签到状态
        await fetchCheckinStats();
        // 更新用户额度
        if (userState.user) {
          const updatedUser = {
            ...userState.user,
            quota: userState.user.quota + data.total_reward,
          };
          userDispatch({ type: 'login', payload: updatedUser });
        }
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('签到失败'));
    } finally {
      setCheckinLoading(false);
    }
  };

  // 执行签到
  const handleCheckin = async () => {
    if (checkinStats?.checked_in_today) {
      Toast.warning(t('今日已签到'));
      return;
    }

    // 如果启用了滑块验证码，先显示验证弹窗
    if (sliderCaptchaEnabled) {
      setCaptchaModalVisible(true);
      return;
    }

    // 否则直接签到
    doCheckin(null);
  };

  const handleCreemCancel = () => {
    setCreemOpen(false);
    setSelectedCreemProduct(null);
  };

  // 选择预设充值额度
  const selectPresetAmount = (preset) => {
    setTopUpCount(preset.value);
    setSelectedPreset(preset.value);

    // 计算实际支付金额，考虑折扣
    const discount = preset.discount || topupInfo.discount[preset.value] || 1.0;
    const discountedAmount = preset.value * priceRatio * discount;
    setAmount(discountedAmount);
  };

  // 格式化大数字显示
  const formatLargeNumber = (num) => {
    return num.toString();
  };

  // 根据最小充值金额生成预设充值额度选项
  const generatePresetAmounts = (minAmount) => {
    const multipliers = [1, 5, 10, 30, 50, 100, 300, 500];
    return multipliers.map((multiplier) => ({
      value: minAmount * multiplier,
    }));
  };

  return (
    <div className='w-full max-w-7xl mx-auto relative min-h-screen lg:min-h-0 mt-[60px] px-2'>
      {/* 划转模态框 */}
      <TransferModal
        t={t}
        openTransfer={openTransfer}
        transfer={transfer}
        handleTransferCancel={handleTransferCancel}
        userState={userState}
        renderQuota={renderQuota}
        getQuotaPerUnit={getQuotaPerUnit}
        transferAmount={transferAmount}
        setTransferAmount={setTransferAmount}
      />

      {/* 充值确认模态框 */}
      <PaymentConfirmModal
        t={t}
        open={open}
        onlineTopUp={onlineTopUp}
        handleCancel={handleCancel}
        confirmLoading={confirmLoading}
        topUpCount={topUpCount}
        renderQuotaWithAmount={renderQuotaWithAmount}
        amountLoading={amountLoading}
        renderAmount={renderAmount}
        payWay={payWay}
        payMethods={payMethods}
        amountNumber={amount}
        discountRate={topupInfo?.discount?.[topUpCount] || 1.0}
      />

      {/* 充值账单模态框 */}
      <TopupHistoryModal
        visible={openHistory}
        onCancel={handleHistoryCancel}
        t={t}
      />

      {/* 签到结果弹窗 */}
      <CheckinResultModal
        visible={checkinResultVisible}
        onClose={() => setCheckinResultVisible(false)}
        result={checkinResult}
        t={t}
        renderQuota={renderQuota}
      />

      {/* 二维码支付弹窗 */}
      <QRCodePaymentModal
        visible={qrCodeModalVisible}
        onClose={handleQRCodeModalClose}
        onSuccess={handleQRCodePaymentSuccess}
        onRefresh={handleQRCodeRefresh}
        qrCodeUrl={qrCodeData.qrCodeUrl}
        tradeNo={qrCodeData.tradeNo}
        amount={qrCodeData.amount}
        money={qrCodeData.money}
        paymentMethod={qrCodeData.paymentMethod}
        t={t}
      />

      {/* 滑块验证码弹窗 */}
      <Modal
        title={t('安全验证')}
        visible={captchaModalVisible}
        onCancel={() => setCaptchaModalVisible(false)}
        footer={null}
        centered
        width={380}
      >
        <div className='py-4'>
          <SliderCaptcha
            onSuccess={handleSliderCaptchaSuccess}
            disabled={checkinLoading}
          />
        </div>
      </Modal>

      {/* Creem 充值确认模态框 */}
      <Modal
        title={t('确定要充值 $')}
        visible={creemOpen}
        onOk={onlineCreemTopUp}
        onCancel={handleCreemCancel}
        maskClosable={false}
        size='small'
        centered
        confirmLoading={confirmLoading}
      >
        {selectedCreemProduct && (
          <>
            <p>
              {t('产品名称')}：{selectedCreemProduct.name}
            </p>
            <p>
              {t('价格')}：{selectedCreemProduct.currency === 'EUR' ? '€' : '$'}{selectedCreemProduct.price}
            </p>
            <p>
              {t('充值额度')}：{selectedCreemProduct.quota}
            </p>
            <p>{t('是否确认充值？')}</p>
          </>
        )}
      </Modal>

      {/* 用户信息头部 */}
      <div className='space-y-6'>
        <div className='grid grid-cols-1 lg:grid-cols-12 gap-6'>
          {/* 左侧充值区域 */}
          <div className='lg:col-span-7 space-y-6 w-full'>
            <RechargeCard
              t={t}
              enableOnlineTopUp={enableOnlineTopUp}
              enableStripeTopUp={enableStripeTopUp}
              enableCreemTopUp={enableCreemTopUp}
              enableLinuxDoTopUp={enableLinuxDoTopUp}
              creemProducts={creemProducts}
              creemPreTopUp={creemPreTopUp}
              presetAmounts={presetAmounts}
              selectedPreset={selectedPreset}
              selectPresetAmount={selectPresetAmount}
              formatLargeNumber={formatLargeNumber}
              priceRatio={priceRatio}
              topUpCount={topUpCount}
              minTopUp={minTopUp}
              renderQuotaWithAmount={renderQuotaWithAmount}
              getAmount={getAmount}
              setTopUpCount={setTopUpCount}
              setSelectedPreset={setSelectedPreset}
              renderAmount={renderAmount}
              amountLoading={amountLoading}
              payMethods={payMethods}
              preTopUp={preTopUp}
              paymentLoading={paymentLoading}
              payWay={payWay}
              redemptionCode={redemptionCode}
              setRedemptionCode={setRedemptionCode}
              topUp={topUp}
              isSubmitting={isSubmitting}
              topUpLink={topUpLink}
              openTopUpLink={openTopUpLink}
              userState={userState}
              renderQuota={renderQuota}
              statusLoading={statusLoading}
              topupInfo={topupInfo}
              onOpenHistory={handleOpenHistory}
              checkinStats={checkinStats}
              checkinLoading={checkinLoading}
              onCheckin={handleCheckin}
            />
          </div>

          {/* 右侧信息区域 */}
          <div className='lg:col-span-5'>
            <InvitationCard
              t={t}
              userState={userState}
              renderQuota={renderQuota}
              setOpenTransfer={setOpenTransfer}
              affLink={affLink}
              handleAffLinkClick={handleAffLinkClick}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default TopUp;
