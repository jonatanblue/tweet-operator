package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Tweet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TweetSpec   `json:"spec,omitempty"`
	Status TweetStatus `json:"status,omitempty"`
}

type TweetSpec struct {
	Text string `json:"text,omitempty"`
}

type TweetStatus struct {
	Likes    int64 `json:"likes,omitempty"`
	Retweets int64 `json:"retweets,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TweetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Tweet `json:"items,omitempty"`
}
