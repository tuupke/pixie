
interface Indexable<T = any> {
    [index: number]: T;
    length: number;
    slice(start: number, end?: number): Indexable<T>;
}

// Similar to a trie. Used to store a list of SequenceKeys
export class Trie<K extends string | number | symbol, V = any> {
    #children: Record<K, Trie<K, V>> = {};
    #value: V | undefined;

    public hasKey(key: K): boolean {
        return this.#children[key] !== undefined;
    }

    public addKey(key: Indexable<K>, value: V): Trie<K, V> {
        if (!this.hasKey(key[0])) {
            this.#children[key[0]] = new Trie();
        }

        if (key.length === 1) {
            this.#children[key[0]].#value = value;
            return this;
        }

        this.#children[key[0]].addKey(key.slice(1), value);
        return this;
    }

    public getValue(key: Indexable<K>): V | undefined {
        if (!key) {
            return undefined;
        }

        if (key.length === 0) {
            return this.#value;
        }

        const subTrie: Trie<K, V> = this.#children[key[0]];
        if (subTrie === undefined) {
            return undefined;
        }

        if (key.length > 1) {
            return subTrie.getValue(key.slice(1));
        }

        return subTrie.#value;
    }

    public forEach(callback: (value: V) => void): void {
        for (const key in this.#children) {
            this.#children[key].forEach(callback);
        }

        if (this.#value !== undefined) {
            callback(this.#value);
        }
    }

    public mergeTriesAt(key: Indexable<K>, value: Trie<K, V>): Trie<K, V> {
        if (!this.hasKey(key[0])) {
            this.#children[key[0]] = new Trie();
        }

        if (key.length === 1) {
            this.#children[key[0]] = value;
            return this;
        }

        this.#children[key[0]].setTrieAt(key.slice(1), value);
        return this;
    }
}
