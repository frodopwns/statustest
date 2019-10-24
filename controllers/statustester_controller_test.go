package controllers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	azurev1 "github.com/Azure/statustest/api/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Status tester Controller", func() {

	BeforeEach(func() {
		// Add any setup steps that needs to be executed before each test
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	// Add Tests for OpenAPI validation (or additonal CRD features) specified in
	// your API definition.
	// Avoid adding tests for vanilla CRUD operations because they would
	// test Kubernetes API server, which isn't the goal here.

	Context("Create and Delete", func() {
		It("should create and delete resource group instances", func() {
			resourceName := "t-rg-dev-statustester"

			var err error
			ctx := context.Background()

			// Create the ResourceGroup object and expect the Reconcile to be created
			resourceInstance := &azurev1.StatusTester{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: "default",
				},
			}

			// create rg
			err = K8sClient.Create(ctx, resourceInstance)
			Expect(apierrors.IsInvalid(err)).To(Equal(false))
			Expect(err).NotTo(HaveOccurred())

			resourceNamespacedName := types.NamespacedName{Name: resourceName, Namespace: "default"}

			// verify rg gets submitted
			Eventually(func() bool {
				_ = K8sClient.Get(ctx, resourceNamespacedName, resourceInstance)
				log.Println(resourceInstance.Status)
				return resourceInstance.Status.Provisioned == true
			}, time.Second*10,
			).Should(BeTrue())

			// delete rg
			err = K8sClient.Delete(ctx, resourceInstance)
			Expect(err).NotTo(HaveOccurred())

			// verify rg is gone from kubernetes
			Eventually(func() bool {
				err := K8sClient.Get(ctx, resourceNamespacedName, resourceInstance)
				if err == nil {
					err = fmt.Errorf("")
				}
				return strings.Contains(err.Error(), "not found")
			}, time.Second*30,
			).Should(BeTrue())

		})
	})
})
