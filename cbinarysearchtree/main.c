#include <stdio.h>
#include <stdlib.h>

struct Node;

struct Node {
    int data;
    struct Node* left;
    struct Node* right;
};

struct Node* init(int data);
void destroy_node(struct Node* node);
void iterate(struct Node* node);
struct Node* find_minimum(struct Node* node);
struct Node* insert(struct Node* node, int item);
struct Node* search(struct Node* node, int item);
struct Node* delete(struct Node* node, int item);

int main(int argc, char** argv) {
    printf("%s\n", "hello");

    struct Node* root = init(6);
    insert(root, 8);
    insert(root, 1);
    insert(root, 20);
    insert(root, 3);
    insert(root, 9);
    insert(root, 2);
    insert(root, 89);
    insert(root, 15);
    insert(root, 5);

    iterate(root);

    printf("--------------\n");

    struct Node* result = search(root, 34);
    printf("search result: %d\n", result->data);

    delete(root, 2);
    printf("--------------\n");
    iterate(root);
    return 0;
}

struct Node* init(int data) {
    struct Node* node = (struct Node*) malloc(sizeof(struct Node*));
    node->data = data;
    node->left = NULL;
    node->right = NULL;
    return node;
}

void destroy_node(struct Node* node) {
    if (node != NULL)
        free((void*) node);
}

void iterate(struct Node* node) {
    if (node == NULL)
        return;
    
    iterate(node->left);
    printf("data: %d\n", node->data);
    iterate(node->right);
}

struct Node* find_minimum(struct Node* node) {
    struct Node* current = node;
    while (current->left != NULL) {
        current = current->left;
    }
    
    return current;
}

struct Node* insert(struct Node* node, int data) {
    if (node == NULL) 
        return init(data);
    
    if (data < node->data) {
        node->left = insert(node->left, data);
    } else {
        node->right = insert(node->right, data);
    }

    return node;
}

struct Node* search(struct Node* node, int data) {
    if (node == NULL) {
        printf("here.. 1\n");
        return insert(node, data);
    }

    if (node->data == data) {
        printf("here.. 2\n");
        return node;
    }
    
    if (data < node->data) {
        return search(node->left, data);
    }

    return search(node->right, data);
}

struct Node* delete(struct Node* node, int data) {
    //base case
    if (node == NULL)
        return NULL;
    
    // if data smaller than the node's data
    // go to the left
    if (data < node->data) {
        node->left = delete(node->left, data);

    // if data greater than the node's data
    // go to the right
    } else if (data > node->data) {
        node->right = delete(node->right, data);
    
    // if data equas to the node's data
    } else {

        // node has no child
        if (node->left == NULL && node->right == NULL) {
            return NULL;
        }

        // node only has one child or has no child
        if (node->left == NULL) {
            struct Node* temp = node->right;
            destroy_node(node);
            return temp;
        }

        if (node->right == NULL) {
            struct Node* temp = node->left;
            destroy_node(node);
            return temp;
        }

        // node has two child
        // find minimum node on the right
        struct Node* min_right_node = find_minimum(node->right);

        // copy data to the successor's data
        node->data = min_right_node->data;

        delete(node->right, min_right_node->data);
    }

    return node;
}