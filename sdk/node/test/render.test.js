import test from 'node:test'
import assert from 'node:assert/strict'
import { renderTemplate, PromptOpsClient } from '../src/index.js'

test('renderTemplate substitutes known variables', () => {
  assert.equal(renderTemplate('hello {{name}}', { name: 'Ada' }), 'hello Ada')
})

test('renderTemplate leaves unknown variables untouched', () => {
  assert.equal(renderTemplate('hello {{name}}', {}), 'hello {{name}}')
})

test('renderTemplate handles whitespace and dotted keys', () => {
  assert.equal(renderTemplate('{{ a.b }}', { 'a.b': 'X' }), 'X')
})

test('renderTemplate substitutes repeated occurrences', () => {
  assert.equal(renderTemplate('{{x}}-{{x}}', { x: 1 }), '1-1')
})

test('constructor requires a server', () => {
  assert.throws(() => new PromptOpsClient({}), /server/)
})

test('constructor strips trailing slashes and applies defaults', () => {
  const c = new PromptOpsClient({ server: 'http://localhost:8080//' })
  assert.equal(c.server, 'http://localhost:8080')
  assert.equal(c.namespace, 'prod')
})
