# Copyright (c) 2015 Platform9 Systems Inc.
# This file provides a simple Levenshtein distance implementation
# on tokens rather than strings

__author__ = 'roopak'


def token_ratio(s1, s2):
    """
    Token ratio returns the similarity index between the two strings based upon their 'token' distance
    tokens are separated by whitespaces
    :return: (ratio, token_len_1, token_len_2)
    """
    tokens_s1 = s1.split()
    tokens_s2 = s2.split()
    largest_len = max(len(tokens_s1), len(tokens_s2))
    dist_metric, subst_tokens, canonical_string = _minimumEditDistance(tokens_s1, tokens_s2)
    return (100 * (largest_len - dist_metric )) / largest_len, len(tokens_s1), len(
            tokens_s2), subst_tokens, canonical_string


def _minimumEditDistance(tokens_s1, tokens_s2):
    """
    EXPERIMENTAL ..

    :param tokens_s1: First
    :param tokens_s2:
    :return:
    """
    insert_cost = 2
    del_cost = 2

    if len(tokens_s1) > len(tokens_s2):
        tokens_s1, tokens_s2 = tokens_s2, tokens_s1
    subst_tokens = []
    canonical_string = []
    last_row = range(len(tokens_s1) + 1)
    last_subst_cost = 0
    last_cost = 0
    for index2, char2 in enumerate(tokens_s2):
        new_row = [index2 + 1]
        min_subst_cost = 1000
        min_subst_val = ""
        min_cost = 1000
        for index1, char1 in enumerate(tokens_s1):
            repel_cost = 1
            if char1 == char2:
                repel_cost = 0
            insertions = new_row[-1] + insert_cost
            deletions = last_row[index1 + 1] + del_cost
            substitutions = last_row[index1] + repel_cost

            min_dist = min(insertions, deletions, substitutions)
            min_cost = min(min_dist, min_cost)
            if repel_cost > 0 and min_dist == substitutions and min_dist < min_subst_cost:
                min_subst_cost = min_dist
                min_subst_val = char1

            new_row.append(min_dist)

        if min_cost == last_cost:
            canonical_string.append(char2)
        else:
            canonical_string.append("%s")
            last_cost = min_cost

        if min_cost == min_subst_cost and len(min_subst_val) > 0 and min_subst_cost > last_subst_cost:
            subst_tokens.append(min_subst_val)
            last_subst_cost = min_subst_cost
        last_row = new_row
    return (last_row[-1], subst_tokens, canonical_string)


def _minimumEditDistance1(tokens_s1, tokens_s2):
    """
    Mimimum distance edit code from rosettacode

    :param tokens_s1: First
    :param tokens_s2:
    :return:
    """
    if len(tokens_s1) > len(tokens_s2):
        tokens_s1, tokens_s2 = tokens_s2, tokens_s1

    distances = range(len(tokens_s1) + 1)
    for index2, char2 in enumerate(tokens_s2):
        newDistances = [index2 + 1]
        for index1, char1 in enumerate(tokens_s1):
            if char1 == char2:
                newDistances.append(distances[index1])
            else:
                min_dist = min((distances[index1],
                                distances[index1 + 1],
                                newDistances[-1]))
                newDistances.append(1 + min_dist)
        distances = newDistances
    return distances[-1]


def _test():
    f = open('./token_match_test.txt')
    data = []
    i = 0
    for l in f.readlines():
        l = l.strip()
        data.append(l)
        if i % 2 != 0:
            print token_ratio(data[0], data[1])
            data = []
        i = i + 1


if __name__ == "__main__":
    s1 = '[resmgr] [Thread-2] Marking the host 87987c71-90f7-4d56-acca-c11adae5c2ed as not responding'
    s2 = '[resmgr] [Thread-3] Marking the host d44f268e-2a7c-488e-a6bf-9d76a0fa6905 as not responding'
    print token_ratio(s1, s2)
    _test()
