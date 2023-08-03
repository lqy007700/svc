package service

import (
	"context"
	"errors"
	"github.com/zxnlx/common"
	"github.com/zxnlx/svc/domain/model"
	"github.com/zxnlx/svc/domain/repository"
	"github.com/zxnlx/svc/proto/svc"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strconv"
)

type ISvcDataService interface {
	AddSvc(svc *model.Svc) (int64, error)
	DelSvc(id int64) error
	UpdateSvc(svc *model.Svc) error
	FindSvcById(id int64) (*model.Svc, error)
	FindAllSvc() ([]model.Svc, error)

	CreateSvcToK8s(svc *svc.SvcInfo) error
	DeleteFromK8s(svc *model.Svc) error
	UpdateSvcToK8s(*svc.SvcInfo) error
}

type SvcDataService struct {
	SvcRepository repository.ISvcRepository
	K8sClientSet  *kubernetes.Clientset
}

func NewSvcDataServices(svcRepo repository.ISvcRepository, clientSet *kubernetes.Clientset) ISvcDataService {
	return &SvcDataService{
		SvcRepository: svcRepo,
		K8sClientSet:  clientSet,
	}
}

func (s *SvcDataService) AddSvc(svc *model.Svc) (int64, error) {
	return s.SvcRepository.CreateSvc(svc)
}

func (s *SvcDataService) DelSvc(id int64) error {
	return s.SvcRepository.DeleteSvcById(id)
}

func (s *SvcDataService) UpdateSvc(svc *model.Svc) error {
	return s.SvcRepository.Update(svc)
}

func (s *SvcDataService) FindSvcById(id int64) (*model.Svc, error) {
	return s.SvcRepository.FindSvcById(id)
}

func (s *SvcDataService) FindAllSvc() ([]model.Svc, error) {
	return s.SvcRepository.FindAll()
}

func (s *SvcDataService) CreateSvcToK8s(svc *svc.SvcInfo) error {
	service := s.setService(svc)
	_, err := s.K8sClientSet.CoreV1().Services(svc.SvcNamespace).Get(context.Background(), svc.SvcName, metav1.GetOptions{})
	if err == nil {
		common.Error("Service " + svc.SvcName + "已经存在")
		return errors.New("Service " + svc.SvcName + "已经存在")
	}

	_, err = s.K8sClientSet.CoreV1().Services(svc.SvcNamespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		common.Error(err)
		return err
	}
	return nil
}

func (s *SvcDataService) DeleteFromK8s(svc *model.Svc) error {
	err := s.K8sClientSet.CoreV1().Services(svc.SvcNamespace).Delete(context.TODO(), svc.SvcName, metav1.DeleteOptions{})
	if err != nil {
		common.Error(err)
		return err
	}
	if err := s.DelSvc(svc.ID); err != nil {
		common.Error(err)
		return err
	}
	common.Info("删除Service ID：" + strconv.FormatInt(svc.ID, 10) + "成功！")
	return nil
}

func (s *SvcDataService) UpdateSvcToK8s(svc *svc.SvcInfo) error {
	service := s.setService(svc)
	//查找是否纯在指定的服务

	_, err := s.K8sClientSet.CoreV1().Services(svc.SvcNamespace).Get(context.TODO(), svc.SvcName, metav1.GetOptions{})
	if err != nil {
		//查找不到
		common.Error(err)
		return errors.New("Service" + svc.SvcName + "不存在请先创建")
	}

	_, err = s.K8sClientSet.CoreV1().Services(svc.SvcNamespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	if err != nil {
		common.Error(err)
		return err
	}
	common.Info("Service " + svc.SvcName + "更新成功")
	return nil
}

func (s *SvcDataService) setService(svc *svc.SvcInfo) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "v1",
			APIVersion: "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svc.SvcName,
			Namespace: svc.SvcNamespace,
			Labels: map[string]string{
				"app-name": svc.SvcPodName,
				"author":   "Caplost",
			},
			Annotations: map[string]string{
				"k8s/generated-by-zxnl": "zxnl",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: s.getSvcPort(svc),
			Selector: map[string]string{
				"app-name": svc.SvcPodName,
			},
			Type: "ClusterIP",
		},
		Status: corev1.ServiceStatus{
			LoadBalancer: corev1.LoadBalancerStatus{},
			Conditions:   nil,
		},
	}
}

func (s *SvcDataService) getSvcPort(svc *svc.SvcInfo) []corev1.ServicePort {
	servicePort := make([]corev1.ServicePort, 0, len(svc.SvcPort))
	for _, v := range svc.SvcPort {
		servicePort = append(servicePort, corev1.ServicePort{
			Name:       "port-" + strconv.FormatInt(int64(v.SvcPort), 10),
			Protocol:   corev1.Protocol(v.SvcPortProtocol),
			Port:       v.SvcPort,
			TargetPort: intstr.FromInt(int(v.SvcTargetPort)),
		})
	}
	return servicePort
}
